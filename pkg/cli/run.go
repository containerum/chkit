package cli

import (
	"github.com/containerum/chkit/pkg/cli/deployment"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/cli/solution"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/spf13/cobra"
)

func Run(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:   "run",
		Short: "Run solutions and deployments",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := prerun.PreRun(ctx); err != nil {
				angel.Angel(ctx, err)
				ctx.Exit(1)
			}
			if err := prerun.GetNamespaceByUserfriendlyID(ctx, cmd.Flags()); err != nil {
				ferr.Println(err)
				ctx.Exit(1)
			}
		},
		Run: func(command *cobra.Command, args []string) {
			command.Help()
		},
		PersistentPostRun: ctx.CobraPostRun,
	}
	command.AddCommand(
		clisolution.Run(ctx),
		clideployment.RunVersion(ctx),
	)
	command.PersistentFlags().
		StringP("namespace", "n", ctx.GetNamespace().ID, "")
	command.PersistentFlags().
		BoolP("help", "h", false, "Print help for chkit")
	return command
}
