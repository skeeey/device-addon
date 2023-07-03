package config

type ProtocolProperties map[string]any

// 'R' 'W' 'RW' 'WR' are supported
type ReadWrite string

type Device struct {
	Name         string                        `yaml:"name"`
	Manufacturer string                        `yaml:"manufacturer"`
	Model        string                        `yaml:"model"`
	Description  string                        `yaml:"description,omitempty"`
	Protocols    map[string]ProtocolProperties `yaml:"protocols"`
	Properties   map[string]any                `yaml:"properties,omitempty"`
	Profile      DeviceProfile                 `yaml:"profile"`
}

type DeviceList struct {
	Devices []Device `yaml:"devices"`
}

type DeviceProfile struct {
	DeviceResources []DeviceResource `yaml:"deviceResources"`
	DeviceCommands  []DeviceCommand  `yaml:"deviceCommands"`
}

type DeviceProfileList struct {
	Profiles []DeviceProfile `yaml:"profiles"`
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
	Name        string             `yaml:"name"`
	Description string             `yaml:"description"`
	Properties  ResourceProperties `yaml:"properties"`
	Attributes  map[string]any     `yaml:"attributes"`
}

type DeviceCommand struct {
	Name      string     `yaml:"name"`
	ReadWrite ReadWrite  `yaml:"readWrite"`
	Resources []Resource `yaml:"resources"`
}

type Resource struct {
	DeviceResource string `yaml:"deviceResource"`
	DefaultValue   any    `yaml:"defaultValue"`
}
