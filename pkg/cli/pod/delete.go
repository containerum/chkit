package clipod

import (
	"strings"

	"os"

	"fmt"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Delete(ctx *context.Context) *cobra.Command {
	var deletePodConfig = struct {
		Force bool
	}{}
	command := &cobra.Command{
		Use:     "pod",
		Aliases: aliases,
		Short:   "delete pod in specific namespace",
		Long:    "deletes pods. Aliases: " + strings.Join(aliases, ", "),
		Example: "chkit delete pod pod_name [-n namespace]",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Help()
				return
			}
			podName := args[0]
			logrus.
				WithField("command", "delete pod").
				Debugf("start deleting pod %q", podName)
			if deletePodConfig.Force || activekit.YesNo(fmt.Sprintf("Are you sure you want to delete pod %q? [Y/N]: ", podName)) {
				if err := ctx.Client.DeletePod(ctx.Namespace.ID, podName); err != nil {
					logrus.WithError(err).Debugf("unable to delete pod %q in namespace %q", podName, ctx.Namespace)
					activekit.Attention(err.Error())
					os.Exit(1)
				}
				fmt.Printf("OK\n")
				logrus.Debugf("pod %q in namespace %q deleted", podName, ctx.Namespace)
			}
		},
	}
	command.PersistentFlags().
		BoolVarP(&deletePodConfig.Force, "force", "f", false, "delete pod without confirmation")
	return command
}
