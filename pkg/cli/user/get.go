package cliuser

import (
	"os"

	"fmt"

	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/spf13/cobra"
)

var aliases = []string{"me", "user"}

func Get(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:     "profile",
		Aliases: aliases,
		Short:   "show profile info",
		Long:    "Shows profile info",
		Example: "chkit get profile",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := prerun.PreRun(ctx); err != nil {
				angel.Angel(ctx, err)
				os.Exit(1)
			}
			if cmd.Flags().Changed("namespace") {
				ctx.Namespace.ID, _ = cmd.Flags().GetString("namespace")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger := coblog.Logger(cmd)
			logger.Debugf("getting profile info")
			profile, err := ctx.Client.GetProfile()
			if err != nil {
				logger.WithError(err).Errorf("unable to get profile info")
				activekit.Attention("Unable to get profile info")
				os.Exit(1)
			}
			fmt.Println(profile)
		},
	}
	return command
}
