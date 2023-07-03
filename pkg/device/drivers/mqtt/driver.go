package mqtt

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"k8s.io/klog/v2"

	"github.com/skeeey/device-addon/pkg/device/config"
	"github.com/skeeey/device-addon/pkg/device/messagebuses"
	"github.com/skeeey/device-addon/pkg/device/models"
	"github.com/skeeey/device-addon/pkg/device/util"
)

type MQTTDriver struct {
	mqttBrokerInfo *MQTTBrokerInfo
	mqttClient     mqtt.Client
	devices        map[string]config.Device
	msgBuses       []messagebuses.MessageBus
}

func NewMQTTDriver() *MQTTDriver {
	return &MQTTDriver{
		devices: make(map[string]config.Device),
	}
}

func (d *MQTTDriver) Initialize(driverInfo config.DriverInfo, msgBuses []messagebuses.MessageBus) error {
	var mqttBrokerInfo = &MQTTBrokerInfo{}
	if err := util.LoadConfig(path.Join(driverInfo.ConfigDir, config.DriverConfigFileName), mqttBrokerInfo); err != nil {
		return err
	}

	if err := d.createMQTTClient(); err != nil {
		return err
	}

	d.msgBuses = msgBuses
	d.mqttBrokerInfo = mqttBrokerInfo
	return nil
}

func (d *MQTTDriver) Start() error {
	//do nothing
	return nil
}

func (d *MQTTDriver) Stop() error {
	klog.Info("driver is stopping, disconnect the MQTT conn")
	if d.mqttClient.IsConnected() {
		d.mqttClient.Disconnect(5000)
	}
	return nil
}

func (d *MQTTDriver) AddDevice(device config.Device) error {
	_, ok := d.devices[device.Name]
	if !ok {
		d.devices[device.Name] = device
	}

	return nil
}

func (d *MQTTDriver) UpdateDevice(device config.Device) error {
	//TODO
	return nil
}

func (d *MQTTDriver) RemoveDevice(deviceName string) error {
	//TODO
	return nil
}

func (d *MQTTDriver) HandleCommands(deviceName string, command models.Command) error {
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

	data := make(models.Attributes)
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
			msgBus.Publish(deviceName, *result)
		}
	}
}
