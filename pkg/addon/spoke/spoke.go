package spoke

import (
	"context"

	"github.com/spf13/pflag"

	"k8s.io/client-go/rest"
)

// AgentOptions defines the flags for workload agent
type AgentOptions struct {
	HubKubeconfigFile string
	SpokeClusterName  string
}

// NewAgentOptions returns the flags with default value set
func NewAgentOptions() *AgentOptions {
	return &AgentOptions{}
}

func (o *AgentOptions) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.HubKubeconfigFile, "hub-kubeconfig", o.HubKubeconfigFile, "Location of kubeconfig file to connect to hub cluster.")
	flags.StringVar(&o.SpokeClusterName, "cluster-name", o.SpokeClusterName, "Name of spoke cluster.")
}

// RunAgent starts the controllers on agent to process work from hub.
func (o *AgentOptions) RunAgent(ctx context.Context, kubeconfig *rest.Config) error {
	//TODO
	return nil
}
