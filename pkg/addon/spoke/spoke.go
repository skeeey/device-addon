package spoke

import (
	"context"
	"time"

	"github.com/spf13/pflag"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/skeeey/device-addon/pkg/addon/spoke/controllers"
	"github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	deviceclient "github.com/skeeey/device-addon/pkg/client/clientset/versioned"
	deviceinformers "github.com/skeeey/device-addon/pkg/client/informers/externalversions"
	"github.com/skeeey/device-addon/pkg/device/equipment"
)

// AgentOptions defines the flags for workload agent
type AgentOptions struct {
	HubKubeconfigFile string
	SpokeClusterName  string
	// TODO read these from addon configuration api
	ReceiveTopic  string
	PayloadFormat string
}

// NewAgentOptions returns the flags with default value set
func NewAgentOptions() *AgentOptions {
	return &AgentOptions{
		ReceiveTopic:  "devices/%s/data/%s",
		PayloadFormat: "jsonMap",
	}
}

func (o *AgentOptions) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.HubKubeconfigFile, "hub-kubeconfig", o.HubKubeconfigFile, "Location of kubeconfig file to connect to hub cluster.")
	flags.StringVar(&o.SpokeClusterName, "cluster-name", o.SpokeClusterName, "Name of spoke cluster.")
	flags.StringVar(&o.ReceiveTopic, "messagebus-receive-topic", o.ReceiveTopic, "")
	flags.StringVar(&o.PayloadFormat, "messagebus-payload-format", o.PayloadFormat, "")
}

// RunAgent starts the controllers on agent to process work from hub.
func (o *AgentOptions) RunAgent(ctx context.Context, kubeconfig *rest.Config) error {
	hubRestConfig, err := clientcmd.BuildConfigFromFlags("", o.HubKubeconfigFile)
	if err != nil {
		return err
	}

	deviceClient, err := deviceclient.NewForConfig(hubRestConfig)
	if err != nil {
		return err
	}

	equipment := equipment.NewEquipment()
	if err := equipment.Start(
		ctx,
		[]v1alpha1.MessageBusConfig{
			{
				MessageBusType: "mqtt",
				Enabled:        true,
				Properties: v1alpha1.Values{
					Data: map[string]interface{}{
						"receiveTopic":  o.ReceiveTopic,
						"payloadFormat": o.PayloadFormat,
					},
				},
			},
		}); err != nil {
		return err
	}

	deviceinformerFactory := deviceinformers.NewSharedInformerFactory(deviceClient, 10*time.Minute)

	driverController := controllers.NewDriversConroller(
		o.SpokeClusterName,
		deviceClient,
		deviceinformerFactory.Edge().V1alpha1().Drivers(),
		equipment,
	)

	deviceController := controllers.NewDevicesConroller(
		o.SpokeClusterName,
		deviceClient,
		deviceinformerFactory.Edge().V1alpha1().Devices(),
		equipment,
	)

	go deviceinformerFactory.Start(ctx.Done())

	go deviceController.Run(ctx, 1)
	go driverController.Run(ctx, 1)

	<-ctx.Done()

	return nil
}
