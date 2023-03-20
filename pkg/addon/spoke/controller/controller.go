package controller

import (
	"context"

	"open-cluster-management.io/addon-framework/pkg/basecontroller/factory"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	deviceclientset "github.com/skeeey/device-addon/pkg/client/clientset/versioned"
	devicesv1alpha1informer "github.com/skeeey/device-addon/pkg/client/informers/externalversions/apis/v1alpha1"
	devicesv1alpha1lister "github.com/skeeey/device-addon/pkg/client/listers/apis/v1alpha1"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
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
				metaObj, _ := meta.Accessor(obj)
				return []string{metaObj.GetName()}
			}, deviceInformer.Informer()).
		WithSync(c.sync).ToController("device-controller")
}

func (c *deviceController) sync(ctx context.Context, syncCtx factory.SyncContext, key string) error {
	return nil
}
