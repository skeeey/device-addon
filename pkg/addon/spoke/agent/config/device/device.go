package device

type ProtocolProperties map[string]any

type Device struct {
	Name        string                        `yaml:"name"`
	Description string                        `yaml:"description,omitempty"`
	Protocols   map[string]ProtocolProperties `yaml:"protocols"`
	Properties  map[string]any                `yaml:"properties,omitempty"`
}

type DeviceList struct {
	Devices []Device `yaml:"devices"`
}
