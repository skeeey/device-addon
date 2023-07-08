// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "github.com/skeeey/device-addon/pkg/client/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Devices returns a DeviceInformer.
	Devices() DeviceInformer
	// DeviceAddOnConfigs returns a DeviceAddOnConfigInformer.
	DeviceAddOnConfigs() DeviceAddOnConfigInformer
	// Drivers returns a DriverInformer.
	Drivers() DriverInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Devices returns a DeviceInformer.
func (v *version) Devices() DeviceInformer {
	return &deviceInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// DeviceAddOnConfigs returns a DeviceAddOnConfigInformer.
func (v *version) DeviceAddOnConfigs() DeviceAddOnConfigInformer {
	return &deviceAddOnConfigInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// Drivers returns a DriverInformer.
func (v *version) Drivers() DriverInformer {
	return &driverInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
