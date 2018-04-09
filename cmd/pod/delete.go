package clipod

import (
	"strings"

	"github.com/containerum/chkit/cmd/cmdutil"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

var Delete = &cli.Command{
	Name:        "pod",
	Usage:       "call to delete pod in specific namespace",
	UsageText:   "chkit delete pod pod_name [-n namespace]",
	Description: "deletes pods. Aliases: " + strings.Join(aliases, ", "),
	Aliases:     aliases,
	Flags:       cmdutil.DeleteFlags,
	Action: func(ctx *cli.Context) error {
		logrus.Debugf("running command delete pod")
		client := cmdutil.GetClient(ctx)
		namespace := cmdutil.GetNamespace(ctx)
		if ctx.NArg() == 0 {
			logrus.Debugf("show help")
			return cli.ShowSubcommandHelp(ctx)
		}
		pod := ctx.Args().First()
		logrus.Debugf("deleting pod %q from %q", pod, namespace)
		err := client.DeletePod(namespace, pod)
		logrus.WithError(err).Debugf("error while deleting pod")
		return err
	},
}
