package set

import (
	"fmt"

	"strings"

	"os"

	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/validation"
	"github.com/spf13/cobra"
)

func DefaultNamespace(ctx *context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "default-namespace",
		Short:   "set default namespace",
		Aliases: []string{"def-ns", "default-ns", "defns", "def-namespace"},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := prerun.PreRun(ctx); err != nil {
				activekit.Attention(err.Error())
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				ctx.Namespace = args[0]
				fmt.Printf("Using %q as default namespace!\n", ctx.Namespace)
				ctx.Changed = true
				return
			}
			nsList, err := ctx.Client.GetNamespaceList()
			if err != nil || len(nsList) == 0 {
				fmt.Printf("You have no namespaces :(\n")
			}
			var menu []*activekit.MenuItem
			for _, ns := range nsList {
				menu = append(menu, &activekit.MenuItem{
					Label: ns.Label,
					Action: func(ns string) func() error {
						return func() error {
							ctx.Namespace = ns
							ctx.Changed = true
							fmt.Printf("Using %q as default namespace\n", ns)
							return nil
						}
					}(ns.ID),
				})
			}
			menu = append(menu, []*activekit.MenuItem{
				{
					Label: "Set custom namespace",
					Action: func() error {
						ns := strings.TrimSpace(activekit.Promt("Type namespace ID: "))
						if err := validation.ValidateID(ns); ns == "" || err != nil {
							fmt.Printf("Inavlid namespace ID\n")
							return nil
						}
						ctx.Namespace = ns
						ctx.Changed = true
						fmt.Printf("Using %q as default namespace!\n", ns)
						return nil
					},
				},
				{
					Label: "Exit",
				},
			}...)
			var title string
			if ctx.Namespace == "" {
				title = fmt.Sprintf("Default namespace isn't defined")
			} else {
				title = fmt.Sprintf("%q is current default namespace", ctx.Namespace)
			}
			(&activekit.Menu{
				Title: title,
				Items: menu,
			}).Run()
		},
	}

}
