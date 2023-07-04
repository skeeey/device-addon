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
	DeviceConfig `json:",inline"`
}

type DeviceStatus struct {
	// conditions describe the state of the current device.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// 'R' 'W' 'RW' 'WR' are supported
type ReadWrite string

type DeviceConfig struct {
	Name         string            `yaml:"name" json:"name"`
	DriverType   string            `yaml:"driverType" json:"driverType"`
	Manufacturer string            `yaml:"manufacturer" json:"manufacturer"`
	Model        string            `yaml:"model" json:"model"`
	Description  string            `yaml:"description,omitempty" json:"description,omitempty"`
	Protocols    map[string]Values `yaml:"protocols" json:"protocols"`
	Properties   Values            `yaml:"properties,omitempty" json:"properties,omitempty"`
	Profile      DeviceProfile     `yaml:"profile" json:"profile"`
}

type DeviceProfile struct {
	DeviceResources []DeviceResource `yaml:"deviceResources" json:"deviceResources"`
	DeviceCommands  []DeviceCommand  `yaml:"deviceCommands" json:"deviceCommands"`
}

type DeviceProfileList struct {
	Profiles []DeviceProfile `yaml:"profiles" json:"profiles"`
}

type ResourceProperties struct {
	ValueType    string    `yaml:"valueType" json:"valueType"`
	ReadWrite    ReadWrite `yaml:"readWrite" json:"readWrite"`
	Units        string    `yaml:"units" json:"units"`
	Minimum      *float64  `yaml:"minimum" json:"minimum"`
	Maximum      *float64  `yaml:"maximum" json:"maximum"`
	DefaultValue string    `yaml:"defaultValue" json:"defaultValue"`
	Mask         *uint64   `yaml:"mask" json:"mask"`
	Shift        *int64    `yaml:"shift" json:"shift"`
	Scale        *float64  `yaml:"scale" json:"scale"`
	Offset       *float64  `yaml:"offset" json:"offset"`
	Base         *float64  `yaml:"base" json:"base"`
	Assertion    string    `yaml:"assertion" json:"assertion"`
	MediaType    string    `yaml:"mediaType" json:"mediaType"`
	Optional     Values    `yaml:"optional" json:"optional"`
}

type DeviceResource struct {
	Name        string             `yaml:"name" json:"name"`
	Description string             `yaml:"description" json:"description"`
	Properties  ResourceProperties `yaml:"properties" json:"properties"`
	Attributes  Values             `yaml:"attributes" json:"attributes"`
}

type DeviceCommand struct {
	Name      string                  `yaml:"name" json:"name"`
	ReadWrite ReadWrite               `yaml:"readWrite" json:"readWrite"`
	Resources []DeviceCommandResource `yaml:"resources" json:"resources"`
}

type DeviceCommandResource struct {
	DeviceResource string `yaml:"deviceResource" json:"deviceResource"`
	DefaultValue   string `yaml:"defaultValue" json:"defaultValue"`
}
