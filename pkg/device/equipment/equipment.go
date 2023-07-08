package equipment

import (
	"fmt"
	"sync"

	"github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	"github.com/skeeey/device-addon/pkg/device/drivers"
	"github.com/skeeey/device-addon/pkg/device/messagebuses"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/klog/v2"
)

type equipmentDriver struct {
	driver drivers.Driver
	config v1alpha1.DriverConfig
}

type Equipment struct {
	sync.Mutex
	messageBuses map[string]messagebuses.MessageBus
	drivers      map[string]equipmentDriver
}

func NewEquipment() *Equipment {
	return &Equipment{
		messageBuses: make(map[string]messagebuses.MessageBus),
		drivers:      make(map[string]equipmentDriver),
	}
}

func (e *Equipment) Start(configs []v1alpha1.MessageBusConfig) error {
	e.Lock()
	defer e.Unlock()

	for _, c := range configs {
		newMsgBus, err := messagebuses.Get(c)
		if newMsgBus == nil {
			continue
		}

		if err != nil {
			return err
		}

		if err := newMsgBus.Start(); err != nil {
			return fmt.Errorf("failed to start message bus %s", c.MessageBusType)
		}

		e.messageBuses[c.MessageBusType] = newMsgBus
	}

	return nil
}

func (e *Equipment) Stop(configs []v1alpha1.MessageBusConfig) {
	e.Lock()
	defer e.Unlock()

	for _, d := range e.drivers {
		d.driver.Stop()
	}

	for _, m := range e.messageBuses {
		m.Stop()
	}
}

func (e *Equipment) InstallDriver(config v1alpha1.DriverConfig) error {
	e.Lock()
	defer e.Unlock()

	msgBuses := []messagebuses.MessageBus{}
	for _, m := range e.messageBuses {
		msgBuses = append(msgBuses, m)
	}

	d := drivers.Get(config.DriverType, config.Properties.Data, msgBuses)
	if d == nil {
		return nil
	}

	if lastDriver, ok := e.drivers[config.DriverType]; ok {
		if equality.Semantic.DeepEqual(lastDriver.config, config) {
			klog.Infof("The driver %s already exists", config.DriverType)
			return nil
		}

		klog.Infof("Reinstall the driver %s already exists", config.DriverType)
		d.Stop()
	}

	if err := d.Start(); err != nil {
		return fmt.Errorf("failed to start driver %s", config.DriverType)
	}

	klog.Infof("The driver %s is installed", config.DriverType)
	e.drivers[config.DriverType] = equipmentDriver{
		driver: d,
		config: config,
	}
	return nil
}

func (e *Equipment) UnInstallDriver(config v1alpha1.DriverConfig) error {
	e.Lock()
	defer e.Unlock()

	d, ok := e.drivers[config.DriverType]
	if !ok {
		klog.Infof("The driver %s does not exist", config.DriverType)
		return nil
	}

	d.driver.Stop()
	delete(e.drivers, config.DriverType)
	return nil
}

func (e *Equipment) GetDriver(driverType string) drivers.Driver {
	e.Lock()
	defer e.Unlock()

	d, ok := e.drivers[driverType]
	if !ok {
		return nil
	}
	return d.driver
}
