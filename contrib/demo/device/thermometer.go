package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	pubTopic = "v1alpha1/devices/attrs/"
	subTopic = "v1alpha1/devices/attrs/push"
)

var temperatures = []int{15, 22, 30}

func init() {
	mqtt.ERROR = log.New(os.Stderr, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stderr, "[CRIT] ", 0)
	// mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	// mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
}

type temperatureHandler func(int)

type Thermometer struct {
	interrupt chan struct{}
	handlers  []temperatureHandler
}

func (t *Thermometer) AddHandlers(handlers ...temperatureHandler) *Thermometer {
	t.handlers = append(t.handlers, handlers...)
	return t
}

func (t *Thermometer) Stop() {
	t.interrupt <- struct{}{}
}

func (t *Thermometer) Run() {
	fmt.Fprintln(os.Stdout, "Starting ...")
	for {
		select {
		case <-t.interrupt:
			fmt.Fprintln(os.Stdout, "Shutdown")
			return
		default:
			for _, handler := range t.handlers {
				handler(temperatures[rand.Intn(len(temperatures)-1)])
			}
			time.Sleep(15 * time.Second)
		}
	}
}

func showTemperature(temperature int) {
	fmt.Fprintln(os.Stdout, "Current temperature: ", temperature)
}

func publish(client mqtt.Client, sensorID string) temperatureHandler {
	return func(temperature int) {
		data := map[string]string{}
		data["temperature"] = strconv.Itoa(temperature)

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to marshal data ", data, ", ", err)
			return
		}

		t := client.Publish(pubTopic+sensorID, 0, false, jsonData)
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

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	defer close(done)

	sensorID := os.Args[1]
	brokeAttr := os.Args[2]

	fmt.Fprintln(os.Stdout, "Thermometer", sensorID, "connecting to MQTT broke:", brokeAttr)
	client, err := connectToMQTT(brokeAttr)
	if err != nil {
		log.Fatal(err)
	}

	thermometer := &Thermometer{
		interrupt: make(chan struct{}),
		handlers:  []temperatureHandler{},
	}
	thermometer.AddHandlers(showTemperature, publish(client, sensorID))

	go thermometer.Run()

	t := client.Subscribe(subTopic+"/+", 0, func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := msg.Payload()

		if strings.Contains(topic, sensorID) {
			attrs := map[string]string{}
			if err := json.Unmarshal(payload, &attrs); err != nil {
				fmt.Fprintln(os.Stderr, "Failed to unmarshal data ", string(payload), ", ", err)
			}

			for k, v := range attrs {
				fmt.Fprintln(os.Stdout, "msg from mqqt broker key=", k, ", value=", v)
			}
		}
	})
	<-t.Done()
	if t.Error() != nil {
		fmt.Fprintln(os.Stderr, "Failed to subscribe mqtt message, ", t.Error())
	}

	<-done

	thermometer.Stop()
}
