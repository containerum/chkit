package clinamespace

import (
	"fmt"

	"os"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/spf13/cobra"
)

var accessAliases = []string{"namespace-access", "ns-access"}

func GetAccess(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:     "access",
		Aliases: accessAliases,
		Short:   "get namespace access",
		Example: "chkit get ns-access $ID",
		Run: func(cmd *cobra.Command, args []string) {
			logger := coblog.Logger(cmd)
			var nsID = ctx.Namespace.ID
			if len(args) == 1 {
				nsID = args[0]
			} else if len(args) > 1 {
				cmd.Help()
				os.Exit(1)
			}
			logger.Debugf("getting namespace %q access", ctx.Namespace)
			acc, err := ctx.Client.GetAccess(nsID)
			if err != nil {
				logger.WithError(err).Errorf("unable to get namespace %q access", nsID)
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(acc.RenderTable())
			return
		},
	}
	return command
}
