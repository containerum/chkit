package cli

import (
	"fmt"

	"github.com/containerum/chkit/pkg/cli/deployment"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/cli/service"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Create = &cobra.Command{
	Use:   "create",
	Short: "Create deployment or service",
	PersistentPreRun: func(command *cobra.Command, args []string) {
		prerun.PreRun()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
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
	Create.AddCommand(
		clideployment.Create,
		cliserv.Create(&context.GlobalContext),
	)
}
