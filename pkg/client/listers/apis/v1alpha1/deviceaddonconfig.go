// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// DeviceAddOnConfigLister helps list DeviceAddOnConfigs.
// All objects returned here must be treated as read-only.
type DeviceAddOnConfigLister interface {
	// List lists all DeviceAddOnConfigs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.DeviceAddOnConfig, err error)
	// DeviceAddOnConfigs returns an object that can list and get DeviceAddOnConfigs.
	DeviceAddOnConfigs(namespace string) DeviceAddOnConfigNamespaceLister
	DeviceAddOnConfigListerExpansion
}

// deviceAddOnConfigLister implements the DeviceAddOnConfigLister interface.
type deviceAddOnConfigLister struct {
	indexer cache.Indexer
}

// NewDeviceAddOnConfigLister returns a new DeviceAddOnConfigLister.
func NewDeviceAddOnConfigLister(indexer cache.Indexer) DeviceAddOnConfigLister {
	return &deviceAddOnConfigLister{indexer: indexer}
}

// List lists all DeviceAddOnConfigs in the indexer.
func (s *deviceAddOnConfigLister) List(selector labels.Selector) (ret []*v1alpha1.DeviceAddOnConfig, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.DeviceAddOnConfig))
	})
	return ret, err
}

// DeviceAddOnConfigs returns an object that can list and get DeviceAddOnConfigs.
func (s *deviceAddOnConfigLister) DeviceAddOnConfigs(namespace string) DeviceAddOnConfigNamespaceLister {
	return deviceAddOnConfigNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// DeviceAddOnConfigNamespaceLister helps list and get DeviceAddOnConfigs.
// All objects returned here must be treated as read-only.
type DeviceAddOnConfigNamespaceLister interface {
	// List lists all DeviceAddOnConfigs in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.DeviceAddOnConfig, err error)
	// Get retrieves the DeviceAddOnConfig from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.DeviceAddOnConfig, error)
	DeviceAddOnConfigNamespaceListerExpansion
}

// deviceAddOnConfigNamespaceLister implements the DeviceAddOnConfigNamespaceLister
// interface.
type deviceAddOnConfigNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all DeviceAddOnConfigs in the indexer for a given namespace.
func (s deviceAddOnConfigNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.DeviceAddOnConfig, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.DeviceAddOnConfig))
	})
	return ret, err
}

// Get retrieves the DeviceAddOnConfig from the indexer for a given namespace and name.
func (s deviceAddOnConfigNamespaceLister) Get(name string) (*v1alpha1.DeviceAddOnConfig, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("deviceaddonconfig"), name)
	}
	return obj.(*v1alpha1.DeviceAddOnConfig), nil
}
