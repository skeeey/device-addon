package drivers

import (
	"context"

	"github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	"github.com/skeeey/device-addon/pkg/device/drivers/mqtt"
	"github.com/skeeey/device-addon/pkg/device/drivers/opcua"
	"github.com/skeeey/device-addon/pkg/device/messagebuses"
	"github.com/skeeey/device-addon/pkg/device/util"
	"k8s.io/klog/v2"
)

type Driver interface {
	Start(ctx context.Context) error

	Stop(ctx context.Context)

	AddDevice(device v1alpha1.DeviceConfig) error

	RemoveDevice(deviceName string) error

	RunCommand(command util.Command) error

	GetType() string
}

func Get(driverType string, driverConfig map[string]interface{}, msgBuses []messagebuses.MessageBus) Driver {
	switch driverType {
	case "mqtt":
		return mqtt.NewMQTTDriver(driverConfig, msgBuses)
	case "opcua":
		return opcua.NewOPCUADriver(driverConfig, msgBuses)
	default:
		klog.Warningf("unsupported driver type %s", driverType)
	}

	return nil
}
