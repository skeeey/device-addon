package models

import "github.com/skeeey/device-addon/pkg/addon/spoke/agent/config/device"

const (
	ValueTypeBool         = "Bool"
	ValueTypeString       = "String"
	ValueTypeUint8        = "Uint8"
	ValueTypeUint16       = "Uint16"
	ValueTypeUint32       = "Uint32"
	ValueTypeUint64       = "Uint64"
	ValueTypeInt8         = "Int8"
	ValueTypeInt16        = "Int16"
	ValueTypeInt32        = "Int32"
	ValueTypeInt64        = "Int64"
	ValueTypeFloat32      = "Float32"
	ValueTypeFloat64      = "Float64"
	ValueTypeBinary       = "Binary"
	ValueTypeBoolArray    = "BoolArray"
	ValueTypeStringArray  = "StringArray"
	ValueTypeUint8Array   = "Uint8Array"
	ValueTypeUint16Array  = "Uint16Array"
	ValueTypeUint32Array  = "Uint32Array"
	ValueTypeUint64Array  = "Uint64Array"
	ValueTypeInt8Array    = "Int8Array"
	ValueTypeInt16Array   = "Int16Array"
	ValueTypeInt32Array   = "Int32Array"
	ValueTypeInt64Array   = "Int64Array"
	ValueTypeFloat32Array = "Float32Array"
	ValueTypeFloat64Array = "Float64Array"
	ValueTypeObject       = "Object"
)

type Attributes map[string]interface{}

type Command struct {
	// refer to DeviceProfile.DeviceCommand.Name
	DeviceCommand string                               `yaml:"DeviceCommand"`
	Protocols     map[string]device.ProtocolProperties `yaml:"protocols"`
	Attributes    Attributes                           `yaml:"atrributes"`
}

type Device struct {
	*device.Device
	*device.DeviceProfile
}

type Result struct {
	Name            string      `json:"name"`
	Value           interface{} `json:"value"`
	Type            string      `json:"type"`
	CreateTimestamp int64       `json:"createTimestamp"`
}
