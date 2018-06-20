package cli

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/pkg/cli/namespace"
	"github.com/containerum/chkit/pkg/cli/postrun"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/spf13/cobra"
)

func Rename(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:   "rename",
		Short: "Rename resource",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := prerun.PreRun(ctx); err != nil {
				angel.Angel(ctx, err)
				os.Exit(1)
			}
			if err := prerun.GetNamespaceByUserfriendlyID(ctx, cmd.Flags()); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PersistentPostRun: postrun.PostRunFunc(ctx),
	}
	command.AddCommand(
		clinamespace.Rename(ctx),
	)
	return command
}
