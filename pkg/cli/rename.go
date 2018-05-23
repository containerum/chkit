package cli

import (
	"github.com/containerum/chkit/pkg/cli/namespace"
	"github.com/containerum/chkit/pkg/cli/postrun"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/spf13/cobra"
)

func Rename(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:   "rename",
		Short: "Rename resource",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := prerun.PreRun(ctx); err != nil {
				angel.Angel(ctx, err)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			postrun.PostRun(coblog.Logger(cmd), ctx)
		},
	}
	command.AddCommand(
		clinamespace.Rename(ctx),
	)
	return command
}
