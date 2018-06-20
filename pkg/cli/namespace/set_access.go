package clinamespace

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func SetAccess(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:        "access",
		Aliases:    accessAliases,
		SuggestFor: accessAliases,
		Short:      "Set namespace access rights",
		Long: "Set namespace access rights.\n" +
			"Available access levels are:\n" + func() string {
			var levels []string
			for _, lvl := range model.Levels() {
				levels = append(levels, lvl.String())
			}
			sort.Strings(levels)
			var lvlsInfo = strings.Join(levels, "\n")
			return text.Indent(lvlsInfo, 2)
		}(),
		Example: "chkit set access $USERNAME $ACCESS_LEVEL [--namespace $ID]",
		PreRun:  prerun.PreRunFunc(ctx),
		Run: func(cmd *cobra.Command, args []string) {
			var logger = coblog.Logger(cmd)
			if len(args) != 2 {
				cmd.Help()
				os.Exit(1)
			}
			var username = args[0]
			accessLevel := model.AccessLevel(args[1])
			if force, _ := cmd.Flags().GetBool("force"); force ||
				activekit.YesNo("Are you sure you want give %s %v access to %s?", username, accessLevel, ctx.Namespace) {
				if err := ctx.Client.SetAccess(ctx.Namespace.ID, username, accessLevel); err != nil {
					logger.WithError(err).Errorf("unable to update access to %q for user %q", username, accessLevel)
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
