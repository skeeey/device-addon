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
	publishTopic  = "publishTopic"
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

func NewMQTTMsgBus() *MQTTMsgBus {
	opts := paho.NewClientOptions()
	opts.AddBroker("tcp://127.0.0.1:1883")
	opts.SetClientID("mqtt-msgbus")
	opts.SetKeepAlive(time.Second * time.Duration(3600))
	opts.SetAutoReconnect(true)

	server := mochi.New(nil)
	l := server.Log.Level(zerolog.ErrorLevel)
	server.Log = &l

	// Allow all connections.
	_ = server.AddHook(new(auth.AllowHook), nil)

	return &MQTTMsgBus{
		mqttBroker: server,
		mqttClient: paho.NewClient(opts),
	}
}

func (m *MQTTMsgBus) Init(msgBusInfo v1alpha1.MessageBusConfig) error {
	ptopic, ok := msgBusInfo.Properties.Data[publishTopic]
	if !ok {
		return fmt.Errorf("the publishTopic is required for mqtt message bus")
	}

	m.pubTopic = fmt.Sprintf("%s", ptopic)

	format, ok := msgBusInfo.Properties.Data[payloadFormat]
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
	//TODO need a embedded flag
	go func() {
		tcp := listeners.NewTCP("mqttmsgbus", ":1883", nil)
		if err := m.mqttBroker.AddListener(tcp); err != nil {
			klog.Fatal(err)
		}

		if err := m.mqttBroker.Serve(); err != nil {
			klog.Fatal(err)
		}

		klog.Infof("The MQTT message bus is started on the local host")
	}()

	// TODO need a notify mechanism
	time.Sleep(5 * time.Second)

	token := m.mqttClient.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	klog.Infof("Connect MQTT message bus")
	return nil
}

func (m *MQTTMsgBus) Publish(deviceName string, result util.Result) {
	topic := fmt.Sprintf(m.pubTopic, deviceName, result.Name)
	data := m.payload(result)

	klog.V(4).Infof("publish the data %s to %s for device %s", string(data), topic, deviceName)
	m.mqttClient.Publish(topic, 0, false, data)
}

func (m *MQTTMsgBus) Subscribe() {
	// TODO subscribe message from message bus and send the message to driver
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
