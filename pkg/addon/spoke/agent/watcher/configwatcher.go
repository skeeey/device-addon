package wacher

import (
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/config"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/config/device"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/drivers"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/models"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/util"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/klog/v2"
)

type DeviceConfigWatcher struct {
	watcher *fsnotify.Watcher
	devices map[string]models.Device
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

	var devices = &device.DeviceList{}
	var profiles = &device.DeviceProfileList{}

	if err := util.LoadConfig(path.Join(configDir, config.DevicesFileName), devices); err != nil {
		return nil, err
	}

	if err := util.LoadConfig(path.Join(configDir, config.DeviceProfilesFileName), profiles); err != nil {
		return nil, err
	}

	deviceWatcher := &DeviceConfigWatcher{
		watcher: watcher,
		devices: make(map[string]models.Device),
		driver:  driver,
	}

	if deviceWatcher.update(devices.Devices, profiles.Profiles) {
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
func (dcw *DeviceConfigWatcher) update(devices []device.Device, profiles []device.DeviceProfile) bool {
	updated := false
	for _, deviceInfo := range devices {
		profile := dcw.findProfile(deviceInfo.Name, profiles)
		if profile == nil {
			klog.Warningf("the device %s does not have profile, ignore", deviceInfo.Name)
			continue
		}

		device := models.Device{
			Device:        &deviceInfo,
			DeviceProfile: profile,
		}

		lastDevice, ok := dcw.devices[deviceInfo.Name]
		if !ok {
			dcw.devices[deviceInfo.Name] = device
			updated = true
			continue
		}

		if equality.Semantic.DeepEqual(lastDevice, device) {
			continue
		}

		dcw.devices[deviceInfo.Name] = device
		updated = true
	}

	return updated
}

func (dcw *DeviceConfigWatcher) findProfile(name string, profiles []device.DeviceProfile) *device.DeviceProfile {
	for _, profile := range profiles {
		if name == profile.DeviceName {
			return &profile
		}
	}

	return nil
}
