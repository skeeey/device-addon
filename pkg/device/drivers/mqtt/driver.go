package mqtt

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"k8s.io/klog/v2"

	"github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	"github.com/skeeey/device-addon/pkg/device/messagebuses"
	"github.com/skeeey/device-addon/pkg/device/util"
)

type MQTTDriver struct {
	mqttBrokerInfo *MQTTBrokerInfo
	mqttClient     mqtt.Client
	devices        map[string]v1alpha1.DeviceConfig
	msgBuses       []messagebuses.MessageBus
}

func NewMQTTDriver(driverConfig util.ConfigProperties, msgBuses []messagebuses.MessageBus) *MQTTDriver {
	var mqttBrokerInfo = &MQTTBrokerInfo{}
	if err := util.ToConfigObj(driverConfig, mqttBrokerInfo); err != nil {
		klog.Errorf("failed to parse mqtt drirver config, %v", err)
		return nil
	}

	return &MQTTDriver{
		devices:        make(map[string]v1alpha1.DeviceConfig),
		msgBuses:       msgBuses,
		mqttBrokerInfo: mqttBrokerInfo,
	}
}

func (d *MQTTDriver) GetType() string {
	return "mqtt"
}

func (d *MQTTDriver) Start() error {
	if err := d.createMQTTClient(); err != nil {
		return err
	}
	return nil
}

func (d *MQTTDriver) Stop() {
	klog.Info("driver is stopping, disconnect the MQTT conn")
	if d.mqttClient.IsConnected() {
		d.mqttClient.Disconnect(5000)
	}
}

func (d *MQTTDriver) AddDevice(device v1alpha1.DeviceConfig) error {
	_, ok := d.devices[device.Name]
	if !ok {
		d.devices[device.Name] = device
	}

	return nil
}

func (d *MQTTDriver) RemoveDevice(deviceName string) error {
	//TODO
	return nil
}

func (d *MQTTDriver) RunCommand(command util.Command) error {
	// TODO
	return nil
}

func (d *MQTTDriver) createMQTTClient() error {
	var client mqtt.Client
	var err error
	for i := 0; i <= d.mqttBrokerInfo.ConnEstablishingRetry; i++ {
		client, err = d.getMQTTClient()
		if err != nil {
			if i >= d.mqttBrokerInfo.ConnEstablishingRetry {
				return err
			}

			klog.Warningf("Unable to connect to MQTT broker, %s, retrying", err)
			time.Sleep(time.Duration(30 * time.Second))
			continue
		}

		break
	}

	d.mqttClient = client
	return nil
}

func (d *MQTTDriver) getMQTTClient() (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()

	opts.AddBroker(fmt.Sprintf("tcp://%s", d.mqttBrokerInfo.Host))
	opts.SetClientID(d.mqttBrokerInfo.ClientId)

	//TODO set username and passwork wtih authMode

	opts.SetKeepAlive(time.Second * time.Duration(d.mqttBrokerInfo.KeepAlive))
	opts.SetAutoReconnect(true)
	opts.OnConnect = d.onConnectHandler

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return client, token.Error()
	}

	klog.Infof("Connect MQTT broke %s with client id %s ", d.mqttBrokerInfo.Host, d.mqttBrokerInfo.ClientId)
	return client, nil
}

func (d *MQTTDriver) onConnectHandler(client mqtt.Client) {
	qos := byte(d.mqttBrokerInfo.Qos)
	incomingTopic := d.mqttBrokerInfo.SubTopic

	token := client.Subscribe(incomingTopic, qos, d.onIncomingDataReceived)
	if token.Wait() && token.Error() != nil {
		client.Disconnect(0)
		klog.Errorf("could not subscribe to topic '%s': %v", incomingTopic, token.Error().Error())
		return
	}
	klog.Infof("Subscribed to topic '%s' for receiving the async reading", incomingTopic)
}

func (d *MQTTDriver) onIncomingDataReceived(_ mqtt.Client, message mqtt.Message) {
	incomingTopic := message.Topic()
	subscribedTopic := d.mqttBrokerInfo.SubTopic
	subscribedTopic = strings.Replace(subscribedTopic, "#", "", -1)
	deviceName := strings.Replace(incomingTopic, subscribedTopic, "", -1)

	device, ok := d.devices[deviceName]
	if !ok {
		klog.Infof("Ignore the unadded device %s", deviceName)
		return
	}

	data := make(util.Attributes)
	if err := json.Unmarshal(message.Payload(), &data); err != nil {
		klog.Errorf("failed to unmarshaling incoming data for device %s, %v", deviceName, err)
		return
	}

	for key, val := range data {
		res := util.FindDeviceResource(key, device.Profile.DeviceResources)
		if res == nil {
			klog.Warningf("The device  %s attribute %s  is unsupported", deviceName, key)
			continue
		}

		result, err := util.NewResult(*res, val)
		if err != nil {
			klog.Errorf("The device %s attribute %s  is unsupported, %v", deviceName, key, err)
			continue
		}

		// publish the message to message bus
		for _, msgBus := range d.msgBuses {
			msgBus.ReceiveData(deviceName, *result)
		}
	}
}
