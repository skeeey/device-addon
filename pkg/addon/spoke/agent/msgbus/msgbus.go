package msgbus

import (
	"fmt"

	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/config"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/models"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/msgbus/mqtt"
)

type MessageBus interface {
	Init(msgBusInfo config.MessageBusInfo) error
	Connect() error
	Publish(deviceName string, result models.Result)
	Subscribe()
}

func Get(msgBusType string, msgBusInfo config.MessageBusInfo) (MessageBus, error) {
	switch msgBusType {
	case "mqtt":
		if !msgBusInfo.Enabled {
			return nil, nil
		}

		mqttMsgBus := mqtt.NewMQTTMsgBus()

		if err := mqttMsgBus.Init(msgBusInfo); err != nil {
			return nil, err
		}

		if err := mqttMsgBus.Connect(); err != nil {
			return nil, err
		}

		return mqttMsgBus, nil
	}

	return nil, fmt.Errorf("unsupported message bus type %s", msgBusType)

}
