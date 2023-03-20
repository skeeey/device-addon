package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope="Cluster"

// DeviceDataModel specifies the device data schema.
type DeviceDataModel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	// spec holds the data schema of a device.
	// +kubebuilder:validation:Required
	// +required
	Spec DeviceDataModelSpec `json:"spec"`
}

// DeviceDataModelList is a list of DeviceDataModel
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DeviceDataModelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []DeviceDataModel `json:"items"`
}

type DeviceDataModelSpec struct {
	// Attributes of a device
	// +optional
	Attributes []Attribute `json:"attributes,omitempty"`
}

// AttributeType represents the stored type of IntOrString.
// +kubebuilder:validation:Enum=int;float;double;bool;string;bytes
type AttributeType string

const (
	Int    AttributeType = "int"
	Float  AttributeType = "float"
	Double AttributeType = "double"
	Bool   AttributeType = "bool"
	String AttributeType = "string"
	Bytes  AttributeType = "bytes"
)

// Attribute describes an individual device attribute.
type Attribute struct {
	// Name of this attribute.
	// +kubebuilder:validation:Required
	// +required
	Name string `json:"name"`

	// Description of this attribute.
	// +optional
	Description string `json:"description,omitempty"`

	// Type of this attribute.
	// +kubebuilder:validation:Required
	// +required
	Type AttributeType `json:"type"`

	// Unit of this attribute
	// +optional
	Unit string `json:"unit,omitempty"`
}
