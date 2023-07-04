package drivers

import (
	"fmt"

	"github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	"github.com/skeeey/device-addon/pkg/device/drivers/mqtt"
	"github.com/skeeey/device-addon/pkg/device/drivers/opcua"
	"github.com/skeeey/device-addon/pkg/device/messagebuses"
	"github.com/skeeey/device-addon/pkg/device/util"
)

type Driver interface {
	Initialize(driverConfig util.ConfigProperties, msgBuses []messagebuses.MessageBus) error

	Start() error

	Stop() error

	AddDevice(device v1alpha1.DeviceConfig) error

	UpdateDevice(device v1alpha1.DeviceConfig) error

	RemoveDevice(deviceName string) error

	HandleCommands(deviceName string, command util.Command) error

	GetType() string
}

func Get(driverType string, driverConfig map[string]interface{}, msgBuses []messagebuses.MessageBus) (Driver, error) {
	switch driverType {
	case "mqtt":
		d := mqtt.NewMQTTDriver()
		err := d.Initialize(driverConfig, msgBuses)
		return d, err
	case "opcua":
		d := opcua.NewOPCUADriver()
		err := d.Initialize(driverConfig, msgBuses)
		return d, err
	}

	return nil, fmt.Errorf("unspported driver type %s", driverType)
}
