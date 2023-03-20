package agent

import (
	"github.com/skeeey/device-addon/pkg/addon/spoke"
	"github.com/spf13/cobra"
	"k8s.io/component-base/version"
	"open-cluster-management.io/addon-framework/pkg/cmd/factory"
)

func NewAgentCommand() *cobra.Command {
	o := spoke.NewAgentOptions()
	cmd := factory.NewControllerCommandConfig("device-addon-agent", version.Get(), o.RunAgent).NewCommand()
	cmd.Use = "agent"
	cmd.Short = "Start the addon agent"

	o.AddFlags(cmd)

	return cmd
}
