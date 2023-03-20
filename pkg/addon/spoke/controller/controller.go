package controller

import (
	"context"

	"open-cluster-management.io/addon-framework/pkg/basecontroller/factory"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	deviceclientset "github.com/skeeey/device-addon/pkg/client/clientset/versioned"
	devicesv1alpha1informer "github.com/skeeey/device-addon/pkg/client/informers/externalversions/apis/v1alpha1"
	devicesv1alpha1lister "github.com/skeeey/device-addon/pkg/client/listers/apis/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

type deviceController struct {
	deviceClient deviceclientset.Interface
	deviceLister devicesv1alpha1lister.DeviceLister
	mqttClient   mqtt.Client
	clusterName  string
}

func NewDeviceController(
	deviceClient deviceclientset.Interface,
	deviceInformer devicesv1alpha1informer.DeviceInformer,
	mqttClient mqtt.Client,
	clusterName string,
) factory.Controller {
	c := &deviceController{
		deviceClient: deviceClient,
		deviceLister: deviceInformer.Lister(),
		mqttClient:   mqttClient,
		clusterName:  clusterName,
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

	klog.Infof(device.Spec.DeviceDataModelRef.Name)

	deviceDataModel, err := c.deviceClient.EdgeV1alpha1().DeviceDataModels().Get(ctx, device.Spec.DeviceDataModelRef.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil
	}

	if err != nil {
		return err
	}

	klog.Infof("%v", deviceDataModel.Spec.Attributes)

	return nil
}
