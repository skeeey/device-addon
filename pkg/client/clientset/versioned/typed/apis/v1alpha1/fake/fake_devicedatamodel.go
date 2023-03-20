// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeDeviceDataModels implements DeviceDataModelInterface
type FakeDeviceDataModels struct {
	Fake *FakeEdgeV1alpha1
}

var devicedatamodelsResource = schema.GroupVersionResource{Group: "edge.open-cluster-management.io", Version: "v1alpha1", Resource: "devicedatamodels"}

var devicedatamodelsKind = schema.GroupVersionKind{Group: "edge.open-cluster-management.io", Version: "v1alpha1", Kind: "DeviceDataModel"}

// Get takes name of the deviceDataModel, and returns the corresponding deviceDataModel object, and an error if there is any.
func (c *FakeDeviceDataModels) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.DeviceDataModel, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(devicedatamodelsResource, name), &v1alpha1.DeviceDataModel{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.DeviceDataModel), err
}

// List takes label and field selectors, and returns the list of DeviceDataModels that match those selectors.
func (c *FakeDeviceDataModels) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.DeviceDataModelList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(devicedatamodelsResource, devicedatamodelsKind, opts), &v1alpha1.DeviceDataModelList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.DeviceDataModelList{ListMeta: obj.(*v1alpha1.DeviceDataModelList).ListMeta}
	for _, item := range obj.(*v1alpha1.DeviceDataModelList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested deviceDataModels.
func (c *FakeDeviceDataModels) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(devicedatamodelsResource, opts))
}

// Create takes the representation of a deviceDataModel and creates it.  Returns the server's representation of the deviceDataModel, and an error, if there is any.
func (c *FakeDeviceDataModels) Create(ctx context.Context, deviceDataModel *v1alpha1.DeviceDataModel, opts v1.CreateOptions) (result *v1alpha1.DeviceDataModel, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(devicedatamodelsResource, deviceDataModel), &v1alpha1.DeviceDataModel{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.DeviceDataModel), err
}

// Update takes the representation of a deviceDataModel and updates it. Returns the server's representation of the deviceDataModel, and an error, if there is any.
func (c *FakeDeviceDataModels) Update(ctx context.Context, deviceDataModel *v1alpha1.DeviceDataModel, opts v1.UpdateOptions) (result *v1alpha1.DeviceDataModel, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(devicedatamodelsResource, deviceDataModel), &v1alpha1.DeviceDataModel{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.DeviceDataModel), err
}

// Delete takes name of the deviceDataModel and deletes it. Returns an error if one occurs.
func (c *FakeDeviceDataModels) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(devicedatamodelsResource, name, opts), &v1alpha1.DeviceDataModel{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeDeviceDataModels) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(devicedatamodelsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.DeviceDataModelList{})
	return err
}

// Patch applies the patch and returns the patched deviceDataModel.
func (c *FakeDeviceDataModels) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.DeviceDataModel, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(devicedatamodelsResource, name, pt, data, subresources...), &v1alpha1.DeviceDataModel{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.DeviceDataModel), err
}