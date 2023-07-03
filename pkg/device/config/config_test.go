package config

import "testing"

func TestLoad(t *testing.T) {
	config, err := LoadConfig("/Users/liuwei/go/src/github.com/edgexfoundry/device-mqtt-go/cmd/res/configuration.yaml")
	t.Errorf("%+v, %v", config, err)
}
