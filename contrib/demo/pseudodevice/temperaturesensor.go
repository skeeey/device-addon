package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/skeeey/device-addon/contrib/demo/pseudodevice/device"
)

func init() {
	mqtt.ERROR = log.New(os.Stderr, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stderr, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	//mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
}

func showTemperature(temperature int) {
	fmt.Fprintln(os.Stdout, "Current temperature: ", temperature)
}

func publishTemperature(client mqtt.Client) device.TemperatureHandler {
	return func(temperature int) {
		t := client.Publish("/device/temperatures/status", 0, false, []byte(strconv.Itoa(temperature)))
		go func() {
			<-t.Done()
			if t.Error() != nil {
				fmt.Fprintln(os.Stderr, "Failed to publish message to mqtt, ", t.Error())
			}
		}()
	}
}

func connectToMQTT(mqttURL string) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttURL)

	client := mqtt.NewClient(opts)
	t := client.Connect()
	<-t.Done()

	return client, t.Error()
}

func subscribe(client mqtt.Client, thermometer *device.Thermometer) error {
	t := client.Subscribe("/device/temperatures/actions", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Fprintln(os.Stderr, "msg from mqqt", msg.Topic(), string(msg.Payload()))
		cmd := string(msg.Payload())
		if cmd == "on" {
			thermometer.TurnOn()
		}

		if cmd == "off" {
			thermometer.TurnOff()
		}
	})
	<-t.Done()
	return t.Error()
}

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	defer close(done)

	thermometer := device.NewThermometer()

	if len(os.Args) == 1 {
		thermometer.TurnOn()
		thermometer.AddHandlers(showTemperature)
		<-done
		thermometer.TurnOff()
		return
	}

	// tcp://127.0.0.1:1883
	fmt.Fprintln(os.Stdout, "Connect to MQTT: ", os.Args[1])
	client, err := connectToMQTT(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	thermometer.AddHandlers(showTemperature, publishTemperature(client))
	if err := subscribe(client, thermometer); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to subscribe mqtt message from: ", os.Args[1])
	}

	<-done
}
