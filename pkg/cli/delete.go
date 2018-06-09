package cli

import (
	"fmt"

	"os"

	"github.com/containerum/chkit/pkg/cli/configmap"
	"github.com/containerum/chkit/pkg/cli/deployment"
	"github.com/containerum/chkit/pkg/cli/ingress"
	"github.com/containerum/chkit/pkg/cli/namespace"
	"github.com/containerum/chkit/pkg/cli/pod"
	"github.com/containerum/chkit/pkg/cli/postrun"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/cli/service"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
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
			if err := prerun.GetNamespaceByUserfriendlyID(ctx, cmd.Flags()); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
		Run: func(command *cobra.Command, args []string) {
			command.Help()
		},
		PersistentPostRun: postrun.PostRunFunc(ctx),
	}
	command.AddCommand(
		clinamespace.Delete(ctx),
		//	clinamespace.DeleteAccess(ctx),
		cliserv.Delete(ctx),
		clideployment.Delete(ctx),
		clipod.Delete(ctx),
		clingress.Delete(ctx),
		cliconfigmap.Delete(ctx),
	)
	command.PersistentFlags().
		StringP("namespace", "n", ctx.Namespace.ID, "")
	return command
}
