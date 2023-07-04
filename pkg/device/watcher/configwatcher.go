package watcher

import (
	"path"

	"github.com/fsnotify/fsnotify"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/klog/v2"

	"github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	"github.com/skeeey/device-addon/pkg/device/drivers"
	"github.com/skeeey/device-addon/pkg/device/util"
)

const devicesConfigFileName = "devices.yaml"

type Equipment struct {
	driver  drivers.Driver
	devices map[string]v1alpha1.DeviceConfig
}

type DeviceList struct {
	Devices []v1alpha1.DeviceConfig
}

type DeviceConfigWatcher struct {
	watcher    *fsnotify.Watcher
	equipments map[string]Equipment
}

func NewDeviceConfigWatcher(configDir string, drivers []drivers.Driver) (*DeviceConfigWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := watcher.Add(configDir); err != nil {
		return nil, err
	}

	deviceList := &DeviceList{}
	if err := util.LoadConfig(path.Join(configDir, devicesConfigFileName), deviceList); err != nil {
		return nil, err
	}

	deviceWatcher := &DeviceConfigWatcher{
		watcher:    watcher,
		equipments: make(map[string]Equipment),
	}

	for _, driver := range drivers {
		deviceWatcher.equipments[driver.GetType()] = Equipment{
			driver:  driver,
			devices: make(map[string]v1alpha1.DeviceConfig),
		}
	}

	for _, equipment := range deviceWatcher.equipments {
		if err := deviceWatcher.update(equipment, deviceList.Devices); err != nil {
			return nil, err
		}
	}

	klog.Infof("Watching the config dir %q", configDir)

	return deviceWatcher, nil
}

func (dcw *DeviceConfigWatcher) Watch() {
	defer dcw.watcher.Close()

	for {
		select {
		case event, ok := <-dcw.watcher.Events:
			if !ok {
				return
			}

			if event.Has(fsnotify.Write) {
				klog.Infof("modified file: %v", event.Name)
				//TODO
			}
		case err, ok := <-dcw.watcher.Errors:
			if !ok {
				return
			}
			klog.Infof("file watch error: %v", err)
		}
	}
}

// TODO also return deleted device list
func (dcw *DeviceConfigWatcher) update(equipment Equipment, devices []v1alpha1.DeviceConfig) error {
	for _, device := range devices {
		if device.DriverType != equipment.driver.GetType() {
			continue
		}

		lastDevice, ok := equipment.devices[device.Name]
		if !ok {
			if err := equipment.driver.AddDevice(device); err != nil {
				return err
			}

			equipment.devices[device.Name] = device
			continue
		}

		if equality.Semantic.DeepEqual(lastDevice, device) {
			continue
		}

		if err := equipment.driver.UpdateDevice(device); err != nil {
			return err
		}

		equipment.devices[device.Name] = device
	}

	return nil
}
