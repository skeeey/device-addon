package watcher

import (
	"path"

	"github.com/fsnotify/fsnotify"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/klog/v2"

	"github.com/skeeey/device-addon/pkg/device/config"
	"github.com/skeeey/device-addon/pkg/device/drivers"
	"github.com/skeeey/device-addon/pkg/device/util"
)

type DeviceConfigWatcher struct {
	watcher *fsnotify.Watcher
	devices map[string]config.Device
	driver  drivers.Driver
}

func NewDeviceConfigWatcher(configDir string, driver drivers.Driver) (*DeviceConfigWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := watcher.Add(configDir); err != nil {
		return nil, err
	}

	deviceList := &config.DeviceList{}
	if err := util.LoadConfig(path.Join(configDir, config.DevicesFileName), deviceList); err != nil {
		return nil, err
	}

	deviceWatcher := &DeviceConfigWatcher{
		watcher: watcher,
		devices: make(map[string]config.Device),
		driver:  driver,
	}

	if deviceWatcher.update(deviceList.Devices) {
		for _, device := range deviceWatcher.devices {
			driver.AddDevice(device)
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
func (dcw *DeviceConfigWatcher) update(devices []config.Device) bool {
	updated := false
	for _, device := range devices {
		lastDevice, ok := dcw.devices[device.Name]
		if !ok {
			dcw.devices[device.Name] = device
			updated = true
			continue
		}

		if equality.Semantic.DeepEqual(lastDevice, device) {
			continue
		}

		dcw.devices[device.Name] = device
		updated = true
	}

	return updated
}
