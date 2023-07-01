package drivers

import (
	"fmt"

	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/config"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/drivers/mqtt"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/drivers/opcua"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/models"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/msgbus"
)

type Driver interface {
	Initialize(driverInfo config.DriverInfo, msgBuses []msgbus.MessageBus) error

	Start() error

	Stop() error

	AddDevice(device models.Device) error

	UpdateDevice(device models.Device) error

	RemoveDevice(deviceName string) error

	HandleCommands(deviceName string, command models.Command) error
}

func Get(driverType string, driverInfo config.DriverInfo, msgBuses []msgbus.MessageBus) (Driver, error) {
	switch driverType {
	case "mqtt":
		d := mqtt.NewMQTTDriver()
		err := d.Initialize(driverInfo, msgBuses)
		return d, err
	case "opcua":
		d := opcua.NewOPCUADriver()
		err := d.Initialize(driverInfo, msgBuses)
		return d, err
	}

	return nil, fmt.Errorf("unspported driver type %s", driverType)
}
