package set

import (
	"fmt"

	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/ninedraft/boxofstuff/str"
	"github.com/spf13/cobra"
)

func DefaultNamespace(ctx *context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "default-project",
		Short:   "Set default project",
		Aliases: []string{"def-pr", "default-pr", "defpr", "def-project"},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var flagNs, _ = cmd.Flags().GetString("project")
			var ns = str.Vector{flagNs}.Append(args...).FirstNonEmpty()
			if err := prerun.PreRun(ctx, prerun.Config{
				NamespaceSelection: prerun.RunNamespaceSelectionAndPersist,
				Namespace:          ns,
			}); err != nil {
				activekit.Attention(err.Error())
				ctx.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Using %q as default project", ctx.GetNamespace())
		},
	}

}
