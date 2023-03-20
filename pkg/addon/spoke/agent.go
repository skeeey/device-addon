package spoke

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/skeeey/device-addon/pkg/addon/spoke/controller"
	deviceclientset "github.com/skeeey/device-addon/pkg/client/clientset/versioned"
	dviceinformer "github.com/skeeey/device-addon/pkg/client/informers/externalversions"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// AgentOptions defines the flags for workload agent
type AgentOptions struct {
	HubKubeconfigFile string
	SpokeClusterName  string
	MQTTBrokerAddr    string
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
}

// RunAgent starts the controllers on agent to process work from hub.
func (o *AgentOptions) RunAgent(ctx context.Context, kubeconfig *rest.Config) error {
	_, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return err
	}

	hubRestConfig, err := clientcmd.BuildConfigFromFlags("", o.HubKubeconfigFile)
	if err != nil {
		return err
	}

	deviceClient, err := deviceclientset.NewForConfig(hubRestConfig)
	if err != nil {
		return err
	}

	deviceInformer := dviceinformer.NewSharedInformerFactory(deviceClient, 10*time.Minute)

	agent := controller.NewDeviceController(
		deviceClient,
		deviceInformer.Edge().V1alpha1().Devices(),
		nil,
		o.SpokeClusterName,
	)

	go deviceInformer.Start(ctx.Done())
	go agent.Run(ctx, 1)

	//
	<-ctx.Done()
	return nil
}
