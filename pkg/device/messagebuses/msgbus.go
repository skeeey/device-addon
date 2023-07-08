package messagebuses

import (
	"fmt"

	"github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	"github.com/skeeey/device-addon/pkg/device/messagebuses/mqtt"
	"github.com/skeeey/device-addon/pkg/device/util"
	"k8s.io/klog/v2"
)

type MessageBus interface {
	Start() error
	Stop()
	ReceiveData(deviceName string, result util.Result) error
	SendData() error
}

func Get(conifg v1alpha1.MessageBusConfig) (MessageBus, error) {
	switch conifg.MessageBusType {
	case "mqtt":
		if conifg.Enabled {
			return mqtt.NewMQTTMsgBus(conifg), nil
		}
	default:
		return nil, fmt.Errorf("unsupported message bus type %s", conifg.MessageBusType)
	}

	klog.Warningf("Thers is no message bus is found")
	return nil, nil
}
