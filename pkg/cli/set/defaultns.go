package set

import (
	"fmt"

	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/spf13/cobra"
)

func DefaultNamespace(ctx *context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "default-namespace",
		Short:   "Set default namespace",
		Aliases: []string{"def-ns", "default-ns", "defns", "def-namespace"},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var ns, _ = cmd.Flags().GetString("namespace")
			if err := prerun.PreRun(ctx, prerun.Config{
				NamespaceSelection: prerun.RunNamespaceSelectionAndPersist,
				Namespace:          ns,
			}); err != nil {
				activekit.Attention(err.Error())
				ctx.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Using %q as default namespace", ctx.GetNamespace())
		},
	}

}
