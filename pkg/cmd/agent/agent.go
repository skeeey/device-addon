package agent

import (
	"context"

	"github.com/spf13/cobra"

	"k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/logs"
	"k8s.io/component-base/version"
	"k8s.io/klog/v2"

	"github.com/skeeey/device-addon/pkg/addon/spoke"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent"
	"open-cluster-management.io/addon-framework/pkg/cmd/factory"
)

func NewAddOnAgentCommand() *cobra.Command {
	o := spoke.NewAgentOptions()
	cmd := factory.NewControllerCommandConfig("device-addon-agent", version.Get(), o.RunAgent).NewCommand()
	cmd.Use = "agent"
	cmd.Short = "Start the addon agent"

	o.AddFlags(cmd.Flags())

	return cmd
}

func NewDriverAgentCommand() *cobra.Command {
	o := agent.NewDriverAgentOptions()
	ctx := context.TODO()
	cmd := &cobra.Command{
		Use:   "driver",
		Short: "Start the driver agent",
		Run: func(cmd *cobra.Command, args []string) {
			logs.InitLogs()

			shutdownCtx, cancel := context.WithCancel(ctx)
			shutdownHandler := server.SetupSignalHandler()
			go func() {
				defer cancel()
				<-shutdownHandler
			}()

			ctx, terminate := context.WithCancel(shutdownCtx)
			defer terminate()

			if err := o.RunDeviceAgent(ctx); err != nil {
				klog.Fatal(err)
			}
		},
	}

	o.AddFlags(cmd.Flags())

	return cmd
}
