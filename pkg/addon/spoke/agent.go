package spoke

import (
	"context"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"

	"github.com/skeeey/device-addon/pkg/addon/spoke/controller"
	deviceclientset "github.com/skeeey/device-addon/pkg/client/clientset/versioned"
	dviceinformer "github.com/skeeey/device-addon/pkg/client/informers/externalversions"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

// AgentOptions defines the flags for workload agent
type AgentOptions struct {
	HubKubeconfigFile string
	SpokeClusterName  string
	MQTTBrokerAddr    string
	// Topic (v1alpha1/devices/attrs/push/+) used to publish the data.
	//  - addon publish the data to the device by this topic
	//  - device subscribe this topic to get the data
	MQTTPublishTopic string

	// Topic (v1alpha1/devices/attrs/+) used to get the data.
	//  - addon subscribe this topic to get the data from the device
	//  - device publish the data by this topic
	MQTTSubscribeTopic string
}

// NewAgentOptions returns the flags with default value set
func NewAgentOptions() *AgentOptions {
	return &AgentOptions{}
}

func (o *AgentOptions) AddFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringVar(&o.HubKubeconfigFile, "hub-kubeconfig", o.HubKubeconfigFile, "Location of kubeconfig file to connect to hub cluster.")
	flags.StringVar(&o.SpokeClusterName, "cluster-name", o.SpokeClusterName, "Name of spoke cluster.")
	flags.StringVar(&o.MQTTBrokerAddr, "mqtt-broker-addr", o.MQTTBrokerAddr, "Address of MQTT broker.")
	flags.StringVar(&o.MQTTPublishTopic, "mqtt-publish-topic", o.MQTTPublishTopic, "Topic of MQTT publish.")
	flags.StringVar(&o.MQTTSubscribeTopic, "mqtt-subscribe-topic", o.MQTTSubscribeTopic, "Topic of MQTT subscribe.")
}

// RunAgent starts the controllers on agent to process work from hub.
func (o *AgentOptions) RunAgent(ctx context.Context, kubeconfig *rest.Config) error {
	mqttClient, err := o.connectToMQTT()
	if err != nil {
		return fmt.Errorf("failed to connect to mqtt broker %s, %v", o.MQTTBrokerAddr, err)
	}

	klog.Infof("Connected to mqtt broker %s", o.MQTTBrokerAddr)

	hubRestConfig, err := clientcmd.BuildConfigFromFlags("", o.HubKubeconfigFile)
	if err != nil {
		return err
	}

	deviceClient, err := deviceclientset.NewForConfig(hubRestConfig)
	if err != nil {
		return err
	}

	deviceInformer := dviceinformer.NewSharedInformerFactoryWithOptions(
		deviceClient, 10*time.Minute, dviceinformer.WithNamespace(o.SpokeClusterName))

	agent := controller.NewDeviceController(
		deviceClient,
		deviceInformer.Edge().V1alpha1().Devices(),
		mqttClient,
		o.SpokeClusterName,
		o.MQTTPublishTopic,
	)

	subscriber := controller.NewSubscriber(
		deviceClient,
		deviceInformer.Edge().V1alpha1().Devices(),
		mqttClient,
		o.SpokeClusterName,
		o.MQTTSubscribeTopic,
	)

	go deviceInformer.Start(ctx.Done())

	go agent.Run(ctx, 1)
	go subscriber.Run(ctx)

	<-ctx.Done()
	return nil
}

func (o *AgentOptions) connectToMQTT() (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(o.MQTTBrokerAddr)

	client := mqtt.NewClient(opts)
	t := client.Connect()
	<-t.Done()

	return client, t.Error()
}
