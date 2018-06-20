package cli

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/pkg/cli/configmap"
	"github.com/containerum/chkit/pkg/cli/deployment"
	"github.com/containerum/chkit/pkg/cli/ingress"
	"github.com/containerum/chkit/pkg/cli/postrun"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/cli/service"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/spf13/cobra"
)

func Create(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "Create resource (deployment, service...)",
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
	command.PersistentFlags().
		StringP("namespace", "n", ctx.Namespace.ID, "")
	command.AddCommand(
		cliconfigmap.Create(ctx),
		clideployment.Create(ctx),
		cliserv.Create(ctx),
		clingress.Create(ctx),
	)
	return command
}
