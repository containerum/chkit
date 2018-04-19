package cli

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/containerum/chkit/pkg/cli/deployment"
	"github.com/containerum/chkit/pkg/cli/namespace"
	"github.com/containerum/chkit/pkg/cli/pod"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/cli/service"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/spf13/cobra"
)

var Get = &cobra.Command{
	Use:   "get",
	Short: "Delete resource data",
	PersistentPreRun: func(command *cobra.Command, args []string) {
		prerun.PreRun()
	},
	Run: func(command *cobra.Command, args []string) {
		command.Help()
	},
	PersistentPostRun: func(command *cobra.Command, args []string) {
		if context.GlobalContext.Changed {
			if err := configuration.SaveConfig(); err != nil {
				logrus.WithError(err).Errorf("unable to save config")
				fmt.Printf("Unable to save config: %v\n", err)
				return
			}
		}
		if err := configuration.SaveTokens(context.GlobalContext.Client.Tokens); err != nil {
			logrus.WithError(err).Errorf("unable to save tokens")
			fmt.Printf("Unable to save tokens: %v\n", err)
			return
		}
	},
}

func init() {
	Get.AddCommand(
		clideployment.Get,
		clinamespace.Get,
		cliserv.Get,
		clipod.Get(&context.GlobalContext),
	)
}
