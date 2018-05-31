package set

import (
	"fmt"

	"os"

	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/namespace"
	"github.com/containerum/chkit/pkg/util/activekit"
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
				ns, err := ctx.Client.GetNamespace(args[0])
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				ctx.SetNamespace(ns)
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
					Label: ns.LabelAndID(),
					Action: func(ns namespace.Namespace) func() error {
						return func() error {
							ctx.SetNamespace(ns)
							fmt.Printf("Using %q as default namespace\n", ns.LabelAndID())
							return nil
						}
					}(ns),
				})
			}
			menu = append(menu, &activekit.MenuItem{
				Label: "Exit",
			})
			var title string
			if ctx.Namespace.IsEmpty() {
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
