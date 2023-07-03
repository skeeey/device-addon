package device

import (
	"context"

	"github.com/spf13/pflag"

	"github.com/skeeey/device-addon/pkg/device/config"
	"github.com/skeeey/device-addon/pkg/device/drivers"
	"github.com/skeeey/device-addon/pkg/device/messagebuses"
	"github.com/skeeey/device-addon/pkg/device/watcher"
)

type DriverAgentOptions struct {
	ConfigFile string
}

func NewDriverAgentOptions() *DriverAgentOptions {
	return &DriverAgentOptions{}
}

func (o *DriverAgentOptions) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.ConfigFile, "driver-config", o.ConfigFile, "Location of driver file")
}

// RunAgent starts the controllers on agent to process work from hub.
func (o *DriverAgentOptions) RunDeviceAgent(ctx context.Context) error {
	config, err := config.LoadConfig(o.ConfigFile)
	if err != nil {
		return err
	}

	msgBuses := []messagebuses.MessageBus{}
	for msgBusType, msgBusInfo := range config.MessageBuses {
		msgBus, err := messagebuses.Get(msgBusType, msgBusInfo)
		if err != nil {
			return err
		}
		if msgBus == nil {
			continue
		}

		msgBuses = append(msgBuses, msgBus)
	}

	for driverType, driverInfo := range config.Drivers {
		driver, err := drivers.Get(driverType, driverInfo, msgBuses)
		if err != nil {
			return err
		}

		configWatcher, err := watcher.NewDeviceConfigWatcher(driverInfo.ConfigDir, driver)
		if err != nil {
			return err
		}

		go configWatcher.Watch()
		go driver.Start()
	}

	<-ctx.Done()
	return nil
}
