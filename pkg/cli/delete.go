package cli

import (
	"fmt"

	"os"

	"github.com/containerum/chkit/pkg/cli/deployment"
	"github.com/containerum/chkit/pkg/cli/ingress"
	"github.com/containerum/chkit/pkg/cli/namespace"
	"github.com/containerum/chkit/pkg/cli/pod"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/cli/service"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Delete(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:   "delete",
		Short: "Delete resource",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := prerun.PreRun(ctx); err != nil {
				angel.Angel(ctx, err)
				os.Exit(1)
			}
			if cmd.Flags().Changed("namespace") {
				ctx.Namespace, _ = cmd.Flags().GetString("namespace")
			}
		},
		Run: func(command *cobra.Command, args []string) {
			command.Help()
		},
		PersistentPostRun: func(command *cobra.Command, args []string) {
			if ctx.Changed {
				if err := configuration.SaveConfig(ctx); err != nil {
					logrus.WithError(err).Errorf("unable to save config")
					fmt.Printf("Unable to save config: %v\n", err)
					return
				}
			}
			if err := configuration.SaveTokens(ctx, ctx.Client.Tokens); err != nil {
				logrus.WithError(err).Errorf("unable to save tokens")
				fmt.Printf("Unable to save tokens: %v\n", err)
				return
			}
		},
	}
	command.AddCommand(
		clinamespace.Delete(ctx),
		cliserv.Delete(ctx),
		clideployment.Delete(ctx),
		clipod.Delete(ctx),
		clingress.Delete(ctx),
	)
	command.PersistentFlags().
		StringP("namespace", "n", ctx.Namespace, "")
	return command
}
