package clinamespace

import (
	"github.com/containerum/chkit/cmd/util"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

var Delete = &cli.Command{
	Name:    "namespace",
	Aliases: aliases,
	Action: func(ctx *cli.Context) error {
		logrus.Debugf("running command delete namespace")
		client := util.GetClient(ctx)
		namespace := util.GetNamespace(ctx)
		if ctx.NArg() == 0 {
			logrus.Debugf("show help")
			return cli.ShowSubcommandHelp(ctx)
		}
		depl := ctx.Args().First()
		logrus.Debugf("deleting deployment %q from %q", depl, namespace)
		err := client.DeleteNamespace(namespace)
		return err
	},
}
