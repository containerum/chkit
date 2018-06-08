package set

import (
	"github.com/containerum/chkit/pkg/cli/containerumapi"
	"github.com/containerum/chkit/pkg/cli/image"
	"github.com/containerum/chkit/pkg/cli/namespace"
	"github.com/containerum/chkit/pkg/cli/postrun"
	"github.com/containerum/chkit/pkg/cli/replicas"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/spf13/cobra"
)

func Set(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:   "set",
		Short: "Set configuration variables",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			postrun.PostRun(coblog.Logger(cmd), ctx)
		},
	}
	command.AddCommand(
		DefaultNamespace(ctx),
		image.Set(ctx),
		replicas.Set(ctx),
		containerumapi.Set(ctx),
		clinamespace.SetAccess(ctx),
	)
	return command
}
