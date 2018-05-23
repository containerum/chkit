package cli

import (
	"os"

	"github.com/containerum/chkit/pkg/cli/configmap"
	"github.com/containerum/chkit/pkg/cli/deployment"
	"github.com/containerum/chkit/pkg/cli/ingress"
	"github.com/containerum/chkit/pkg/cli/postrun"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/cli/service"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/spf13/cobra"
)

func Create(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "Create deployment or service",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := prerun.PreRun(ctx); err != nil {
				angel.Angel(ctx, err)
				os.Exit(1)
			}
			if cmd.Flags().Changed("namespace") {
				ctx.Namespace, _ = cmd.Flags().GetString("namespace")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PersistentPostRun: func(command *cobra.Command, args []string) {
			postrun.PostRun(coblog.Logger(command), ctx)
		},
	}
	command.PersistentFlags().
		StringP("namespace", "n", ctx.Namespace, "")
	command.AddCommand(
		cliconfigmap.Create(ctx),
		clideployment.Create(ctx),
		cliserv.Create(ctx),
		clingress.Create(ctx),
	)
	return command
}
