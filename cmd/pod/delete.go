package clipod

import (
	"github.com/containerum/chkit/cmd/util"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

var Delete = &cli.Command{
	Name:    "pod",
	Aliases: aliases,
	Flags:   util.DeleteFlags,
	Action: func(ctx *cli.Context) error {
		logrus.Debugf("running command delete pod")
		client := util.GetClient(ctx)
		namespace := util.GetNamespace(ctx)
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
