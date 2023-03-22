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
	// deviceDataModelRef refers to a device data model.
	// +kubebuilder:validation:Required
	// +required
	DeviceDataModelRef *DeviceDataModelReference `json:"deviceDataModelRef"`

	// data lists the device data that will be processed by the device.
	// +optional
	Data []DeviceData `json:"data,omitempty"`

	// desiredAttrs lists desired device attributes that will be reported from the device.
	// +optional
	DesiredAttrs []string `json:"desiredAttrs,omitempty"`
}

type DeviceStatus struct {
	// conditions describe the state of the current device.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// reportedAttrs contains desired device attributes that are reported from the device.
	// +optional
	ReportedAttrs []ReportedAttr `json:"reportedAttrs,omitempty"`
}

type DeviceDataModelReference struct {
	// name of the DeviceDataMode referent.
	// +kubebuilder:validation:Required
	// +required
	Name string `json:"name"`
}

type DeviceData struct {
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

	// name of a device attribue
	// +kubebuilder:validation:Required
	// +required
	Name string `json:"name"`

	// value of a device attribue
	// +kubebuilder:validation:Required
	// +optional
	Value string `json:"value"`
}
