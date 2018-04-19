package set

import (
	"fmt"

	"github.com/containerum/chkit/pkg/cli/image"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/containerum/chkit/pkg/cli/replicas"
)

func Set(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:   "set",
		Short: "Set configuration variables",
		PersistentPreRun: func(command *cobra.Command, args []string) {
			prerun.PreRun()
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PersistentPostRun: func(command *cobra.Command, args []string) {
			if ctx.Changed {
				if err := configuration.SaveConfig(); err != nil {
					logrus.WithError(err).Errorf("unable to save config")
					fmt.Printf("Unable to save config: %v\n", err)
					return
				}
			}
			if err := configuration.SaveTokens(ctx.Client.Tokens); err != nil {
				logrus.WithError(err).Errorf("unable to save tokens")
				fmt.Printf("Unable to save tokens: %v\n", err)
				return
			}
		},
	}
	command.AddCommand(
		DefaultNamespace(ctx),
		image.Set(ctx),
		replicas.Set(ctx),
	)
	return command
}
