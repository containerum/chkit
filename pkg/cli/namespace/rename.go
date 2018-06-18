package clinamespace

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/namespace"
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
		Example: "chkit rename ns $ID $NEW_NAME",
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
				(&activekit.Menu{
					Title: "Select namespace to rename",
					Items: activekit.StringSelector(nsList.OwnersAndLabels(), func(label string) error {
						ns, err := prerun.ResolveLabel(ctx, label)
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
						interactiveRename(ctx, logger, ns)
						return nil
					}),
				}).Run()
			case 1:
				ns, err := prerun.ResolveLabel(ctx, args[0])
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				interactiveRename(ctx, logger, ns)
				return
			case 2:
				oldNs, err := prerun.ResolveLabel(ctx, args[0])
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				newName := args[1]
				if err := validation.ValidateLabel(newName); err != nil {
					fmt.Printf("invalid new namespace name:\n%v\n", err)
					os.Exit(1)
				}
				if force, _ := cmd.Flags().GetBool("force"); force ||
					activekit.YesNo("Are you sure you want to rename namespace %q?", oldNs.OwnerAndLabel()) {
					if err := ctx.Client.RenameNamespace(oldNs.ID, newName); err != nil {
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

func interactiveRename(ctx *context.Context, logger logrus.FieldLogger, ns namespace.Namespace) {
	for {
		newName := activekit.Promt("Type new namespace name: ")
		if err := validation.ValidateLabel(newName); err != nil {
			fmt.Printf("invalid new namespace name:\n%v\n", err)
			continue
		}
		if activekit.YesNo("Are you sure you want to rename namespace %q?", ns.OwnerAndLabel()) {
			if err := ctx.Client.RenameNamespace(ns.ID, newName); err != nil {
				logger.WithError(err).Errorf("unable to rename namespace %q", ns.OwnerAndLabel())
				fmt.Println(err)
				continue
			}
			fmt.Println("OK")
		}
		return
	}
}
