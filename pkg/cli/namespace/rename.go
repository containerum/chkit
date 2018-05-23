package clinamespace

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/containerum/chkit/pkg/util/validation"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Rename(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:     "namespace",
		Aliases: aliases,
		Run: func(cmd *cobra.Command, args []string) {
			var logger = coblog.Logger(cmd)
			switch len(args) {
			case 0:
				nsList, err := ctx.Client.GetNamespaceList()
				if err != nil {
					logger.WithError(err).Errorf("unable to get namespace list")
					fmt.Println(err)
					os.Exit(1)
				}
				var menu activekit.MenuItems = make([]*activekit.MenuItem, 0, len(nsList))
				for _, ns := range nsList {
					menu = menu.Append(&activekit.MenuItem{
						Label: ns.Label,
						Action: func(nsName string) func() error {
							return func() error {
								interactiveRename(ctx, logger, nsName)
								return nil
							}
						}(ns.Label),
					})
				}
				(&activekit.Menu{
					Title: "Select namespace to rename",
					Items: menu,
				}).Run()
			case 1:
				interactiveRename(ctx, logger, args[0])
				return
			case 2:
				nsName := args[0]
				newName := args[1]
				if err := validation.ValidateLabel(newName); err != nil {
					fmt.Printf("invalid new namespace name:\n%v\n", err)
					os.Exit(1)
				}
				if force, _ := cmd.Flags().GetBool("force"); force ||
					activekit.YesNo("Are you sure you want to rename namespace %q?", nsName) {
					if err := ctx.Client.RenameNamespace(nsName, newName); err != nil {
						logger.WithError(err).Errorf("unable to rename namespace %q")
						fmt.Println(err)
						os.Exit(1)
					}
					fmt.Println("OK")
				}
				return
			default:
				cmd.Help()
				os.Exit(1)
			}
		},
	}
	return command
}

func interactiveRename(ctx *context.Context, logger logrus.FieldLogger, nsName string) {
	for {
		newName := activekit.Promt("Type new namespace name: ")
		if err := validation.ValidateLabel(newName); err != nil {
			fmt.Printf("invalid new namespace name:\n%v\n", err)
			continue
		}
		if activekit.YesNo("Are you sure you want to rename namespace %q?", nsName) {
			if err := ctx.Client.RenameNamespace(nsName, newName); err != nil {
				logger.WithError(err).Errorf("unable to rename namespace %q", nsName)
				fmt.Println(err)
				continue
			}
			fmt.Println("OK")
		}
		return
	}
}
