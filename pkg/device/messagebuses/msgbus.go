package messagebuses

import (
	"fmt"

	"github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	"github.com/skeeey/device-addon/pkg/device/messagebuses/mqtt"
	"github.com/skeeey/device-addon/pkg/device/util"
)

type MessageBus interface {
	Init(msgBusInfo v1alpha1.MessageBusConfig) error
	Connect() error
	Publish(deviceName string, result util.Result)
	Subscribe()
	Stop()
}

func Get(msgBusType string, msgBusInfo v1alpha1.MessageBusConfig) (MessageBus, error) {
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
