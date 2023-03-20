// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// DeviceLister helps list Devices.
// All objects returned here must be treated as read-only.
type DeviceLister interface {
	// List lists all Devices in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Device, err error)
	// Devices returns an object that can list and get Devices.
	Devices(namespace string) DeviceNamespaceLister
	DeviceListerExpansion
}

// deviceLister implements the DeviceLister interface.
type deviceLister struct {
	indexer cache.Indexer
}

// NewDeviceLister returns a new DeviceLister.
func NewDeviceLister(indexer cache.Indexer) DeviceLister {
	return &deviceLister{indexer: indexer}
}

// List lists all Devices in the indexer.
func (s *deviceLister) List(selector labels.Selector) (ret []*v1alpha1.Device, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Device))
	})
	return ret, err
}

// Devices returns an object that can list and get Devices.
func (s *deviceLister) Devices(namespace string) DeviceNamespaceLister {
	return deviceNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// DeviceNamespaceLister helps list and get Devices.
// All objects returned here must be treated as read-only.
type DeviceNamespaceLister interface {
	// List lists all Devices in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Device, err error)
	// Get retrieves the Device from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Device, error)
	DeviceNamespaceListerExpansion
}

// deviceNamespaceLister implements the DeviceNamespaceLister
// interface.
type deviceNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Devices in the indexer for a given namespace.
func (s deviceNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Device, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Device))
	})
	return ret, err
}

// Get retrieves the Device from the indexer for a given namespace and name.
func (s deviceNamespaceLister) Get(name string) (*v1alpha1.Device, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("device"), name)
	}
	return obj.(*v1alpha1.Device), nil
}
