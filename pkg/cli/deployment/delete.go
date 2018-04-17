package clideployment

import (
	"strings"

	"fmt"
	"os"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Delete(ctx *context.Context) *cobra.Command {
	var deleteDeplConfig = struct {
		Force bool
	}{}
	command := &cobra.Command{
		Use:     "deployment",
		Aliases: aliases,
		Short:   "call to delete deployment in specific namespace",
		Long:    "call to delete deployment in specific namespace. Aliases: " + strings.Join(aliases, ", "),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Help()
				return
			}
			deplName := args[0]
			logrus.
				WithField("command", "delete deployment").
				Debugf("start deleting deployment %q", deplName)
			if deleteDeplConfig.Force || activekit.YesNo(fmt.Sprintf("Are you sure you want to delete deployment %q? [Y/N]: ", deplName)) {
				if err := ctx.Client.DeletePod(ctx.Namespace, deplName); err != nil {
					logrus.WithError(err).Debugf("unable to delete deployment %q in namespace %q", deplName, ctx.Namespace)
					activekit.Attention(err.Error())
					os.Exit(1)
				}
				fmt.Printf("OK\n")
				logrus.Debugf("deployment %q in namespace %q deleted", deplName, ctx.Namespace)
			}
		},
	}
	command.PersistentFlags().
		BoolVarP(&deleteDeplConfig.Force, "force", "f", false, "delete without confirmation")
	return command
}
