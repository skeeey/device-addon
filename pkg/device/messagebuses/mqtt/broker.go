package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	mochi "github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/hooks/auth"
	"github.com/mochi-co/mqtt/v2/listeners"
	"github.com/rs/zerolog"

	"k8s.io/klog/v2"

	"github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	"github.com/skeeey/device-addon/pkg/device/util"
)

const (
	publishTopic  = "receiveTopic"
	payloadFormat = "payloadFormat"
)

const (
	jsonObj = "jsonObj"
	jsonMap = "jsonMap"
)

type palyloadFunc func(util.Result) []byte

type MQTTMsgBus struct {
	mqttBroker *mochi.Server
	mqttClient paho.Client
	pubTopic   string
	payload    palyloadFunc
}

func NewMQTTMsgBus(config v1alpha1.MessageBusConfig) *MQTTMsgBus {
	opts := paho.NewClientOptions()
	opts.AddBroker("tcp://127.0.0.1:1883")
	opts.SetClientID("msgbus-mqtt-client")
	opts.SetKeepAlive(time.Second * time.Duration(3600))
	opts.SetAutoReconnect(true)

	server := mochi.New(nil)
	l := server.Log.Level(zerolog.ErrorLevel)
	server.Log = &l

	// Allow all connections.
	_ = server.AddHook(new(auth.AllowHook), nil)

	ptopic, ok := config.Properties.Data[publishTopic]
	if !ok {
		klog.Infof("Using %s as the default publish topic devices/+/data/+", ptopic)
		ptopic = "devices/%s/data/%s"
	}

	format, ok := config.Properties.Data[payloadFormat]
	if !ok {
		klog.Infof("Using %s as the default payload format", jsonMap)
		format = jsonMap
	}

	m := &MQTTMsgBus{
		mqttBroker: server,
		mqttClient: paho.NewClient(opts),
		pubTopic:   fmt.Sprintf("%s", ptopic),
	}

	format = fmt.Sprintf("%s", format)
	switch format {
	case jsonObj:
		m.payload = toJsonObj
	case jsonMap:
		m.payload = toJsonMap
	}

	return m
}

func (m *MQTTMsgBus) Start() error {
	go func() {
		tcp := listeners.NewTCP("mqttmsgbus", ":1883", nil)
		if err := m.mqttBroker.AddListener(tcp); err != nil {
			klog.Fatal(err)
		}

		if err := m.mqttBroker.Serve(); err != nil {
			klog.Fatal(err)
		}

		klog.Infof("MQTT message bus is started on the localhost")
	}()

	// TODO need a notify mechanism
	time.Sleep(5 * time.Second)

	token := m.mqttClient.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	klog.Infof("Connect to localhost MQTT message bus")
	return nil
}

func (m *MQTTMsgBus) ReceiveData(deviceName string, result util.Result) error {
	topic := fmt.Sprintf(m.pubTopic, deviceName, result.Name)
	data := m.payload(result)

	klog.Infof("Send data to MQTT message bus, [%s] [%s] %s", topic, deviceName, string(data))
	m.mqttClient.Publish(topic, 0, false, data)
	return nil
}

func (m *MQTTMsgBus) SendData() error {
	// TODO subscribe to message bus to get the command to send the command to driver
	return nil
}

func (m *MQTTMsgBus) Stop() {
	m.mqttClient.Disconnect(1000)
	m.mqttBroker.Close()
}

func toJsonObj(result util.Result) []byte {
	payload, _ := json.Marshal(result)
	return payload
}

func toJsonMap(result util.Result) []byte {
	payload, _ := json.Marshal(map[string]any{result.Name: result.Value})
	return payload
}
