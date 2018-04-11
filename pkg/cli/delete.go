package cli

import (
	"fmt"

	"github.com/containerum/chkit/pkg/cli/namespace"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/configuration"
	. "github.com/containerum/chkit/pkg/context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Delete = &cobra.Command{
	Use: "delete",
	PersistentPreRun: func(command *cobra.Command, args []string) {
		prerun.PreRun()
	},
	Run: func(command *cobra.Command, args []string) {
		command.Help()
	},
	PersistentPostRun: func(command *cobra.Command, args []string) {
		if Context.Changed {
			if err := configuration.SaveConfig(); err != nil {
				logrus.WithError(err).Errorf("unable to save config")
				fmt.Printf("Unable to save config: %v\n", err)
				return
			}
		}
		if err := configuration.SaveTokens(Context.Client.Tokens); err != nil {
			logrus.WithError(err).Errorf("unable to save tokens")
			fmt.Printf("Unable to save tokens: %v\n", err)
			return
		}
	},
}

func init() {
	Delete.AddCommand(
		clinamespace.Delete,
	)
}
