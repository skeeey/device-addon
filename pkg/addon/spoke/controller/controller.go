package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"open-cluster-management.io/addon-framework/pkg/basecontroller/factory"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	edgev1alpha1 "github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	deviceclientset "github.com/skeeey/device-addon/pkg/client/clientset/versioned"
	devicesv1alpha1informer "github.com/skeeey/device-addon/pkg/client/informers/externalversions/apis/v1alpha1"
	devicesv1alpha1lister "github.com/skeeey/device-addon/pkg/client/listers/apis/v1alpha1"
	"github.com/skeeey/device-addon/pkg/utils"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

const deviceFinalizer = "edge.open-cluster-management.io/api-resource-cleanup"

type deviceController struct {
	deviceClient deviceclientset.Interface
	deviceLister devicesv1alpha1lister.DeviceLister
	mqttClient   mqtt.Client
	// TODO time series data
	devices     map[string][]byte
	clusterName string
	pubTopic    string
}

func NewDeviceController(
	deviceClient deviceclientset.Interface,
	deviceInformer devicesv1alpha1informer.DeviceInformer,
	mqttClient mqtt.Client,
	clusterName string,
	pubTopic string,
) factory.Controller {
	c := &deviceController{
		deviceClient: deviceClient,
		deviceLister: deviceInformer.Lister(),
		mqttClient:   mqttClient,
		devices:      map[string][]byte{},
		clusterName:  clusterName,
		pubTopic:     pubTopic,
	}
	return factory.New().
		WithInformersQueueKeysFunc(
			func(obj runtime.Object) []string {
				key, _ := cache.MetaNamespaceKeyFunc(obj)
				return []string{key}
			}, deviceInformer.Informer()).
		WithSync(c.sync).ToController("device-controller")
}

func (c *deviceController) sync(ctx context.Context, syncCtx factory.SyncContext, key string) error {
	klog.Infof("Reconciling addon deploy %q", key)

	clusterName, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return nil
	}

	device, err := c.deviceLister.Devices(clusterName).Get(name)
	if errors.IsNotFound(err) {
		return nil
	}

	if err != nil {
		return err
	}

	if device.DeletionTimestamp.IsZero() {
		hasFinalizer := false
		for i := range device.Finalizers {
			if device.Finalizers[i] == deviceFinalizer {
				hasFinalizer = true
				break
			}
		}
		if !hasFinalizer {
			finalizerBytes, err := json.Marshal(append(device.Finalizers, deviceFinalizer))
			if err != nil {
				return err
			}
			patch := fmt.Sprintf("{\"metadata\": {\"finalizers\": %s}}", string(finalizerBytes))

			_, err = c.deviceClient.EdgeV1alpha1().Devices(clusterName).Patch(
				ctx, device.Name, types.MergePatchType, []byte(patch), metav1.PatchOptions{})
			return err
		}
	}

	if !device.DeletionTimestamp.IsZero() {
		delete(c.devices, device.Name)
		return c.removeFinalizer(ctx, device)
	}

	modelName := device.Spec.DeviceDataModelRef.Name
	if _, err := c.deviceClient.EdgeV1alpha1().DeviceDataModels().Get(ctx, modelName, metav1.GetOptions{}); err != nil {
		return err
	}

	return c.publishData(ctx, device.DeepCopy())
}

func (c *deviceController) publishData(ctx context.Context, device *edgev1alpha1.Device) error {
	data := device.Spec.Data
	if len(data) == 0 {
		return nil
	}

	dataMap := map[string]any{}
	for _, attr := range data {
		// TODO convert the value with data model type
		dataMap[attr.Name] = attr.Value
	}

	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}

	lastData, ok := c.devices[device.Name]
	if !ok {
		c.devices[device.Name] = jsonData
	}

	if bytes.Equal(lastData, jsonData) {
		return nil
	}

	t := c.mqttClient.Publish(fmt.Sprintf("%s/%s", c.pubTopic, device.Name), 0, false, jsonData)
	go func() {
		<-t.Done()
		publishedCondition := metav1.Condition{
			Type:    "DataPublished",
			Status:  metav1.ConditionTrue,
			Reason:  "Device data is published",
			Message: fmt.Sprintf("Device data %q is published", string(jsonData)),
		}

		// TODO think about how to hanlde errors
		if t.Error() != nil {
			fmt.Fprintln(os.Stderr, "Failed to publish message to mqtt, ", t.Error())
			publishedCondition.Status = metav1.ConditionFalse
			publishedCondition.Reason = "Failed to publish device data"
			publishedCondition.Message = fmt.Sprintf("Failed to publish device data %q, %v", string(jsonData), t.Error())
		}

		if _, _, err := utils.UpdateDeviceStatus(
			ctx,
			c.deviceClient,
			device.Namespace,
			device.Name,
			utils.UpdateDeviceConditionFn(publishedCondition),
		); err != nil {
			klog.Errorf("failed to update device %s/%s status, %v", device.Namespace, device.Name, err)
		}
	}()

	return nil
}

func (c *deviceController) removeFinalizer(ctx context.Context, device *edgev1alpha1.Device) error {
	copiedFinalizers := []string{}
	for i := range device.Finalizers {
		if device.Finalizers[i] == deviceFinalizer {
			continue
		}
		copiedFinalizers = append(copiedFinalizers, device.Finalizers[i])
	}

	if len(device.Finalizers) != len(copiedFinalizers) {
		finalizerBytes, err := json.Marshal(copiedFinalizers)
		if err != nil {
			return err
		}
		patch := fmt.Sprintf("{\"metadata\": {\"finalizers\": %s}}", string(finalizerBytes))

		_, err = c.deviceClient.EdgeV1alpha1().Devices(device.Namespace).Patch(
			ctx, device.Name, types.MergePatchType, []byte(patch), metav1.PatchOptions{})
		return err
	}

	return nil
}
