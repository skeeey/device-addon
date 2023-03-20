package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status

// Device is the Schema for the devices API
type Device struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	// spec holds the information about a device.
	// +kubebuilder:validation:Required
	// +required
	Spec DeviceSpec `json:"spec"`

	// status holds the state of a device.
	// +optional
	Status DeviceStatus `json:"status,omitempty"`
}

// DeviceList is a list of Device
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DeviceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Device `json:"items"`
}

type DeviceSpec struct {
	// id of a device.
	// +kubebuilder:validation:Required
	// +required
	ID string `json:"id"`

	// DeviceDataModelRef refers to a device data model.
	// +kubebuilder:validation:Required
	// +required
	DeviceDataModelRef *DeviceDataModelReference `json:"deviceDataModelRef"`

	// Data will be processed by the device.
	// +optional
	Data DeviceData `json:"data,omitempty"`

	// DesiredData lists desired device attributes that will be reported from the device.
	// +optional
	DesiredData DesiredDeviceData `json:"desiredData,omitempty"`
}

type DeviceStatus struct {
	// conditions describe the state of the current device.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// ReportedAttrs contains desired device attributes that are reported from current device.
	// +optional
	ReportedAttrs []ReportedAttr `json:"reportedAttrs,omitempty"`
}

type DeviceDataModelReference struct {
	// Name of the DeviceDataMode referent.
	// +kubebuilder:validation:Required
	// +required
	Name string `json:"name"`
}

type DeviceData struct {
	// Topic used to publish the data.
	//  - addon publish the data to the device by this topic
	//  - device subscribe this topic to get the data
	// +kubebuilder:default=v1alpha1/devices/+/attrs/push
	// +optional
	Topic string `json:"topic,omitempty"`

	// Attrs lists the data will be published to the device.
	// +optional
	Attrs []AttrData `json:"attrs,omitempty"`
}

type DesiredDeviceData struct {
	// Topic used to get the data.
	//  - addon subscribe this topic to get the data from the device
	//  - device publish the data by this topic
	// +kubebuilder:default=v1alpha1/devices/+/attrs
	// +optional
	Topic string `json:"topic,omitempty"`

	// Attrs lists desired device attribute names that will be reported from current device.
	// +optional
	Attrs []string `json:"attrs,omitempty"`
}

type AttrData struct {
	// name of a device data
	// +kubebuilder:validation:Required
	// +required
	Name string `json:"name"`

	// value of a device data
	// +kubebuilder:validation:Required
	// +required
	Value string `json:"value"`
}

type ReportedAttr struct {
	// lastUpdatedTime is the last updated time for this attribute.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Format=date-time
	// +kubebuilder:validation:Type=string
	// +required
	LastUpdatedTime metav1.Time `json:"lastUpdatedTime"`

	DeviceData DeviceData `json:",inline"`
}
