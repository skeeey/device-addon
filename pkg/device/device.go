package device

import (
	"context"
	"path"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"

	"github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	"github.com/skeeey/device-addon/pkg/device/drivers"
	"github.com/skeeey/device-addon/pkg/device/messagebuses"
	"github.com/skeeey/device-addon/pkg/device/util"
	"github.com/skeeey/device-addon/pkg/device/watcher"
)

const driverConfigFileName = "config.yaml"

type DriverAgentOptions struct {
	ConfigDir string
}

func NewDriverAgentOptions() *DriverAgentOptions {
	return &DriverAgentOptions{}
}

func (o *DriverAgentOptions) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.ConfigDir, "config-dir", o.ConfigDir, "Directory of config files")
}

// RunAgent starts the controllers on agent to process work from hub.
func (o *DriverAgentOptions) RunDriverAgent(ctx context.Context) error {
	config := &v1alpha1.Config{}
	if err := util.LoadConfig(path.Join(o.ConfigDir, driverConfigFileName), config); err != nil {
		return err
	}

	klog.Infof("-----------> %+v", config)

	msgBuses := []messagebuses.MessageBus{}
	for mType, mConfig := range config.MessageBuses {
		msgBus, err := messagebuses.Get(mType, mConfig)
		if err != nil {
			return err
		}
		if msgBus == nil {
			continue
		}

		msgBuses = append(msgBuses, msgBus)
	}

	allDrivers := []drivers.Driver{}
	for driverType, driverConfig := range config.Drivers {
		driver, err := drivers.Get(driverType, driverConfig.Data, msgBuses)
		if err != nil {
			return err
		}

		allDrivers = append(allDrivers, driver)
	}

	configWatcher, err := watcher.NewDeviceConfigWatcher(o.ConfigDir, allDrivers)
	if err != nil {
		return err
	}

	go configWatcher.Watch()

	for _, d := range allDrivers {
		go d.Start()
	}

	<-ctx.Done()
	return nil
}
