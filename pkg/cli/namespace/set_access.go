package clinamespace

import (
	"os"

	"fmt"

	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/access"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func SetAccess(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:     "access",
		Aliases: accessAliases,
		Short:   "get namespace access",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := prerun.PreRun(ctx); err != nil {
				angel.Angel(ctx, err)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			var logger = coblog.Logger(cmd)
			var accessLevel = access.None
			var user string
			var ns string
			switch len(args) {
			case 2:
				ns = args[0]
				user = args[1]
				var menu activekit.MenuItems
				for _, lvl := range access.Levels() {
					menu = menu.Append(&activekit.MenuItem{
						Label: lvl.String(),
						Action: func(lvl access.AccessLevel) func() error {
							return func() error {
								accessLevel = lvl
								return nil
							}
						}(lvl),
					})
				}
				(&activekit.Menu{
					Title: "Select access level",
					Items: menu,
				}).Run()
			case 3:
				ns = args[0]
				user = args[1]
				acc, err := access.LevelFromString(args[2])
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				accessLevel = acc
			default:
				cmd.Help()
				os.Exit(1)
			}
			if force, _ := cmd.Flags().GetBool("force"); force ||
				activekit.YesNo("Are you sure you want give %s %v access?", user, accessLevel) {
				if err := ctx.Client.SetAccess(ns, user, accessLevel); err != nil {
					logger.WithError(err).Errorf("unable to update access to %q for user %q", user, accessLevel)
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println("OK")
			}
		},
	}
	command.PersistentFlags().
		BoolP("force", "f", false, "suppress confirmation")
	return command
}

func selectNamespace(ctx *context.Context, logger logrus.FieldLogger) string {
	nsList, err := ctx.Client.GetNamespaceList()
	if err != nil {
		logger.WithError(err).Errorf("unable to get namespace list")
		fmt.Println(err)
		os.Exit(1)
	}
	var ns string
	var menu activekit.MenuItems
	for _, n := range nsList {
		menu = menu.Append(&activekit.MenuItem{
			Label: n.Label,
			Action: func(nsName string) func() error {
				return func() error {
					ns = nsName
					return nil
				}
			}(n.Label),
		})
	}
	(&activekit.Menu{
		Title: "Select namespace",
		Items: menu,
	}).Run()
	return ns
}
