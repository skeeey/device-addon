package device

// 'R' 'W' 'RW' 'WR'
type ReadWrite string

type DeviceProfile struct {
	DeviceProfileBasicInfo `yaml:",inline"`
	DeviceResources        []DeviceResource `yaml:"deviceResources"`
	DeviceCommands         []DeviceCommand  `yaml:"deviceCommands"`
}

type DeviceProfileList struct {
	Profiles []DeviceProfile `yaml:"profiles"`
}

type DeviceProfileBasicInfo struct {
	DeviceName   string `yaml:"deviceName"`
	Manufacturer string `yaml:"manufacturer"`
	Description  string `yaml:"description"`
	Model        string `yaml:"model"`
}

type ResourceProperties struct {
	ValueType    string         `yaml:"valueType"`
	ReadWrite    ReadWrite      `yaml:"readWrite"`
	Units        string         `yaml:"units"`
	Minimum      *float64       `yaml:"minimum"`
	Maximum      *float64       `yaml:"maximum"`
	DefaultValue string         `yaml:"defaultValue"`
	Mask         *uint64        `yaml:"mask"`
	Shift        *int64         `yaml:"shift"`
	Scale        *float64       `yaml:"scale"`
	Offset       *float64       `yaml:"offset"`
	Base         *float64       `yaml:"base"`
	Assertion    string         `yaml:"assertion"`
	MediaType    string         `yaml:"mediaType"`
	Optional     map[string]any `yaml:"optional"`
}

type DeviceResource struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Properties  ResourceProperties     `yaml:"properties"`
	Attributes  map[string]interface{} `yaml:"attributes"`
	Tags        map[string]any         `yaml:"tags,omitempty"`
}

type DeviceCommand struct {
	Name               string              `yaml:"name"`
	ReadWrite          ReadWrite           `yaml:"readWrite"`
	Tags               map[string]any      `yaml:"tags,omitempty"`
	ResourceOperations []ResourceOperation `yaml:"resourceOperations"`
}

type ResourceOperation struct {
	DeviceResource string `yaml:"deviceResource"`
	DefaultValue   any    `yaml:"defaultValue"`
}
