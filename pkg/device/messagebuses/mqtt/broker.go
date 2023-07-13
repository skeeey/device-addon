package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eclipse/paho.golang/paho"
	mochi "github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/hooks/auth"
	"github.com/mochi-co/mqtt/v2/listeners"
	"github.com/rs/zerolog"

	"k8s.io/klog/v2"

	"github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	"github.com/skeeey/device-addon/pkg/device/client"
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
	pubClient  *paho.Client
	pubTopic   string
	payload    palyloadFunc
}

func NewMQTTMsgBus(config v1alpha1.MessageBusConfig) *MQTTMsgBus {
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

func (m *MQTTMsgBus) Start(ctx context.Context) error {
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

	client, err := client.ConnectToMQTTBroker(
		ctx,
		&client.MQTTBrokerInfo{
			Host:      "127.0.0.1:1883",
			ClientId:  "msgbus-mqtt-pub-client",
			KeepAlive: 3600,
		},
		nil,
	)
	if err != nil {
		return err
	}

	m.pubClient = client

	klog.Infof("Connect to localhost MQTT message bus")
	return nil
}

func (m *MQTTMsgBus) ReceiveData(deviceName string, result util.Result) error {
	topic := fmt.Sprintf(m.pubTopic, deviceName, result.Name)
	data := m.payload(result)

	klog.Infof("Send data to MQTT message bus, [%s] [%s] %s", topic, deviceName, string(data))
	_, err := m.pubClient.Publish(context.TODO(), &paho.Publish{
		Topic:   topic,
		QoS:     0,
		Payload: data,
	})
	if err != nil {
		// TODO handle this error
		klog.Errorf("failed to send data, %v", err)
		return nil
	}

	return nil
}

func (m *MQTTMsgBus) SendData() error {
	// TODO subscribe to message bus to get the command to send the command to driver
	return nil
}

func (m *MQTTMsgBus) Stop(ctx context.Context) {
	m.pubClient.Disconnect(&paho.Disconnect{ReasonCode: 0})
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
