package agent

import (
	"context"

	"github.com/spf13/pflag"

	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/config"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/drivers"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/msgbus"
	wacher "github.com/skeeey/device-addon/pkg/addon/spoke/agent/watcher"
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

	msgBuses := []msgbus.MessageBus{}
	for msgBusType, msgBusInfo := range config.MessageBuses {
		msgBus, err := msgbus.Get(msgBusType, msgBusInfo)
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

		watcher, err := wacher.NewDeviceConfigWatcher(driverInfo.ConfigDir, driver)
		if err != nil {
			return err
		}

		go watcher.Watch()
		go driver.Start()
	}

	<-ctx.Done()
	return nil
}
