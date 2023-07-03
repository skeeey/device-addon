package drivers

import (
	"fmt"

	"github.com/skeeey/device-addon/pkg/device/config"
	"github.com/skeeey/device-addon/pkg/device/drivers/mqtt"
	"github.com/skeeey/device-addon/pkg/device/drivers/opcua"
	"github.com/skeeey/device-addon/pkg/device/messagebuses"
	"github.com/skeeey/device-addon/pkg/device/models"
)

type Driver interface {
	Initialize(driverInfo config.DriverInfo, msgBuses []messagebuses.MessageBus) error

	Start() error

	Stop() error

	AddDevice(device config.Device) error

	UpdateDevice(device config.Device) error

	RemoveDevice(deviceName string) error

	HandleCommands(deviceName string, command models.Command) error
}

func Get(driverType string, driverInfo config.DriverInfo, msgBuses []messagebuses.MessageBus) (Driver, error) {
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
