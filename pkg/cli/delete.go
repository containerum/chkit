package cli

import (
	"github.com/containerum/chkit/pkg/cli/configmap"
	"github.com/containerum/chkit/pkg/cli/deployment"
	"github.com/containerum/chkit/pkg/cli/ingress"
	"github.com/containerum/chkit/pkg/cli/namespace"
	"github.com/containerum/chkit/pkg/cli/pod"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/cli/service"
	"github.com/containerum/chkit/pkg/cli/solution"
	"github.com/containerum/chkit/pkg/cli/volume"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/spf13/cobra"
)

func Delete(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:   "delete",
		Short: "Delete resource",
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
		clinamespace.Delete(ctx),
		//	clinamespace.DeleteAccess(ctx),
		volume.Delete(ctx),
		cliserv.Delete(ctx),
		clideployment.Delete(ctx),
		clideployment.DeleteContainer(ctx),
		clipod.Delete(ctx),
		clingress.Delete(ctx),
		cliconfigmap.Delete(ctx),
		clisolution.Delete(ctx),
	)
	command.PersistentFlags().
		StringP("namespace", "n", ctx.GetNamespace().ID, "")
	return command
}
