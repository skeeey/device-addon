package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

const (
	DevicesFileName        = "devices.yaml"
	DeviceProfilesFileName = "devices.profile.yaml"
	DriverConfigFileName   = "config.yaml"
)

type DriverConfig struct {
	Drivers      map[string]DriverInfo     `yaml:"drivers"`
	MessageBuses map[string]MessageBusInfo `yaml:"messageBuses"`
}

type DriverInfo struct {
	ConfigDir string `yaml:"configDir"`
}

type MessageBusInfo struct {
	Enabled    bool           `yaml:"enabled"`
	Properties map[string]any `yaml:"properties"`
}

func LoadConfig(configFile string) (*DriverConfig, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	config := &DriverConfig{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}
