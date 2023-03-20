// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/skeeey/device-addon/pkg/client/clientset/versioned/typed/apis/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeEdgeV1alpha1 struct {
	*testing.Fake
}

func (c *FakeEdgeV1alpha1) Devices(namespace string) v1alpha1.DeviceInterface {
	return &FakeDevices{c, namespace}
}

func (c *FakeEdgeV1alpha1) DeviceDataModels() v1alpha1.DeviceDataModelInterface {
	return &FakeDeviceDataModels{c}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeEdgeV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
