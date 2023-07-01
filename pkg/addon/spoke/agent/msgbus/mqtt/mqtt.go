package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/config"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/models"
	"k8s.io/klog/v2"
)

const (
	publishTopic  = "publishTopic"
	payloadFormat = "payloadFormat"
)

const (
	jsonObj = "jsonObj"
	jsonMap = "jsonMap"
)

type palyloadFunc func(models.Result) []byte

type MQTTMsgBus struct {
	mqttClient mqtt.Client
	pubTopic   string
	payload    palyloadFunc
}

func NewMQTTMsgBus() *MQTTMsgBus {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://127.0.0.1:1883")
	opts.SetClientID("mqtt-msgbus")
	opts.SetKeepAlive(time.Second * time.Duration(3600))
	opts.SetAutoReconnect(true)

	return &MQTTMsgBus{mqttClient: mqtt.NewClient(opts)}
}

func (m *MQTTMsgBus) Init(msgBusInfo config.MessageBusInfo) error {
	ptopic, ok := msgBusInfo.Properties[publishTopic]
	if !ok {
		return fmt.Errorf("the publishTopic is required for mqtt message bus")
	}

	m.pubTopic = fmt.Sprintf("%s", ptopic)

	format, ok := msgBusInfo.Properties[payloadFormat]
	if !ok {
		return nil
	}

	format = fmt.Sprintf("%s", format)
	switch format {
	case jsonObj:
		m.payload = toJsonObj
	case jsonMap:
		m.payload = toJsonMap
	default:
		return fmt.Errorf("unsupported payload formt for mqtt message bus")
	}
	return nil
}

func (m *MQTTMsgBus) Connect() error {
	token := m.mqttClient.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	klog.Infof("Connect MQTT message bus")
	return nil
}

func (m *MQTTMsgBus) Publish(deviceName string, result models.Result) {
	topic := fmt.Sprintf(m.pubTopic, deviceName, result.Name)
	data := m.payload(result)

	klog.Infof("publish the data %s to %s for device %s", string(data), topic, deviceName)
	m.mqttClient.Publish(topic, 0, false, data)
}

func (m *MQTTMsgBus) Subscribe() {
	// TODO subscribe message from message bus and send the message to driver
}

func toJsonObj(result models.Result) []byte {
	payload, _ := json.Marshal(result)
	return payload
}

func toJsonMap(result models.Result) []byte {
	payload, _ := json.Marshal(map[string]any{result.Name: result.Value})
	return payload
}
