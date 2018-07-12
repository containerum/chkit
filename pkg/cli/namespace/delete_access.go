package clinamespace

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/spf13/cobra"
)

func DeleteAccess(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:        "access",
		Aliases:    accessAliases,
		SuggestFor: accessAliases,
		Short:      "delete user access to namespace",
		Example:    "chkit delete access $USERNAME [--namespace $ID]",
		Run: func(cmd *cobra.Command, args []string) {
			logger := coblog.Logger(cmd)
			if len(args) != 1 {
				cmd.Help()
				ctx.Exit(1)
			}
			username := args[0]
			if force, _ := cmd.Flags().GetBool("force"); force ||
				activekit.YesNo("Are you sure you want to delete %s access to namespace %s?", username, ctx.GetNamespace()) {
				if err := ctx.Client.DeleteAccess(ctx.GetNamespace().ID, username); err != nil {
					logger.WithError(err).Errorf("unable to delete access %s to namespace %s", username, ctx.GetNamespace())
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Println("OK")
			}
		},
	}
	command.PersistentFlags().
		BoolP("force", "f", false, "suppress confirmation")
	return command
}
