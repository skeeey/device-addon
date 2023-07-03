package util

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"

	"github.com/skeeey/device-addon/pkg/device/config"
	"github.com/skeeey/device-addon/pkg/device/models"
)

const castError = "fail to parse %v reading, %v"

func LoadConfig(configFile string, config any) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return err
	}

	return nil
}

func NewResult(resource config.DeviceResource, reading interface{}) (*models.Result, error) {
	var err error
	valueType := resource.Properties.ValueType
	if !checkValueInRange(valueType, reading) {
		return nil, fmt.Errorf("unsupported type %s in device resource", valueType)
	}

	var val interface{}
	switch valueType {
	case models.ValueTypeBool:
		val, err = cast.ToBoolE(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
	case models.ValueTypeString:
		val, err = cast.ToStringE(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
	case models.ValueTypeUint8:
		val, err = cast.ToUint8E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
	case models.ValueTypeUint16:
		val, err = cast.ToUint16E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
	case models.ValueTypeUint32:
		val, err = cast.ToUint32E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
	case models.ValueTypeUint64:
		val, err = cast.ToUint64E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
	case models.ValueTypeInt8:
		val, err = cast.ToInt8E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
	case models.ValueTypeInt16:
		val, err = cast.ToInt16E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
	case models.ValueTypeInt32:
		val, err = cast.ToInt32E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
	case models.ValueTypeInt64:
		val, err = cast.ToInt64E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
	case models.ValueTypeFloat32:
		val, err = cast.ToFloat32E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
	case models.ValueTypeFloat64:
		val, err = cast.ToFloat64E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
	case models.ValueTypeObject:
		val = reading
	default:
		return nil, fmt.Errorf("return result fail, none supported value type: %v", valueType)

	}

	return &models.Result{
		Name:            resource.Name,
		Type:            valueType,
		Value:           val,
		CreateTimestamp: time.Now().UnixNano(),
	}, nil
}

func FindDeviceResource(name string, resources []config.DeviceResource) *config.DeviceResource {
	for _, res := range resources {
		if res.Name == name {
			return &res
		}
	}

	return nil
}
