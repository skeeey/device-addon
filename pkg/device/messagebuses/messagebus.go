package messagebuses

import (
	"context"
	"fmt"

	"github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	"github.com/skeeey/device-addon/pkg/device/messagebuses/mqtt"
	"github.com/skeeey/device-addon/pkg/device/util"

	"k8s.io/klog/v2"
)

type MessageBus interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context)
	ReceiveData(deviceName string, result util.Result) error
	SendData() error
}

func Get(config v1alpha1.MessageBusConfig) (MessageBus, error) {
	switch config.MessageBusType {
	case "mqtt":
		if config.Enabled {
			return mqtt.NewMQTTMsgBus(config), nil
		}
	default:
		return nil, fmt.Errorf("unsupported message bus type %s", config.MessageBusType)
	}

	klog.Warningf("There is no message bus is found")
	return nil, nil
}
