package clinamespace

import (
	"fmt"
	"sync/atomic"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/namespace"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/containerum/chkit/pkg/util/limiter"
	"github.com/ninedraft/boxofstuff/str"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Delete(ctx *context.Context) *cobra.Command {
	var flags struct {
		Force bool     `desc:"suppress confirmation"`
		ID    []string `desc:"project id to delete"`
		Label []string `desc:"project label or owner/label to delete"`
	}
	command := &cobra.Command{
		Use:     "project",
		Short:   "delete project",
		Long:    "Delete project provided in the first arg.",
		Example: "chkit delete project $ID",
		Aliases: aliases,
		Run: func(command *cobra.Command, args []string) {
			var logger = ctx.Log.Command("delete namespace")
			var namespacesToDelete namespace.NamespaceList
			logger.Debugf("getting namespace lst")
			var namespaces, err = ctx.Client.GetNamespaceList()
			if err != nil {
				logger.WithError(err).Debugf("unable to get namespace list")
				ferr.Println(err)
				ctx.Exit(1)
			}
			if len(args) > 0 || len(flags.Label) > 0 {
				var labels = str.Vector(args).Append(flags.Label...).Unique()
				namespacesToDelete = namespaces.Filter(func(i namespace.Namespace) bool {
					return labels.Filter(func(label string) bool { return i.MatchLabel(label) }).Len() > 0
				})
			}
			if len(flags.ID) > 0 {
				namespacesToDelete = namespaces.Filter(func(i namespace.Namespace) bool {
					return str.Vector(flags.ID).Contains(i.ID)
				})
			}
			switch {
			case flags.Force && namespacesToDelete.Len() == 0:
				ferr.Printf("namespaces to delete must be defined as args, --id flag or --label flag")
				ctx.Exit(1)
			case !flags.Force && namespacesToDelete.Len() == 0:
				for exit := false; !exit; {
					var idsToDelete str.Vector = namespacesToDelete.IDs()
					namespaces = namespaces.Filter(func(i namespace.Namespace) bool {
						return !idsToDelete.Contains(i.ID)
					})
					(&activekit.Menu{
						Title: "Select namespace",
						Items: activekit.ItemsFromIter(uint(namespaces.Len()), func(index uint) *activekit.MenuItem {
							var ns = namespaces[index]
							return &activekit.MenuItem{
								Label: ns.OwnerAndLabel(),
								Action: func() error {
									namespacesToDelete = namespacesToDelete.Append(ns)
									idsToDelete = idsToDelete.Append(ns.ID)
									return nil
								},
							}
						}).Append(&activekit.MenuItem{
							Label: "Confirm",
							Action: func() error {
								exit = true
								return nil
							},
						}),
					}).Run()
				}
			}
			if flags.Force || activekit.YesNo("Do you really want to delete namespaces?") {
				if namespacesToDelete.Len() == 0 {
					return
				}
				var limit = limiter.New(4)
				var exitCode uint32
				for _, ns := range namespacesToDelete {
					go func(done func(), ns namespace.Namespace) {
						logger.Debugf("deleting namespace %q", ns.OwnerAndLabel())
						defer done()
						if err := ctx.Client.DeleteNamespace(ns.ID); err != nil {
							logger.WithError(err).Errorf("unable to delete namespace %q", ns.OwnerAndLabel())
							ferr.Printf("unable to delete namespace: %v\n", err)
							atomic.AddUint32(&exitCode, 1)
						}
					}(limit.Start(), ns)
				}
				limit.Wait()
				if exitCode == 0 {
					fmt.Println("Ok")
				}
				ctx.Exit(int(exitCode))
			}
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
