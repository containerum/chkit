package clideployment

import (
	"github.com/containerum/chkit/cmd/util"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

var DeleteDeployment = &cli.Command{
	Name:    "deployment",
	Aliases: []string{"depl", "deployments"},
	Action: func(ctx *cli.Context) error {
		logrus.Debugf("running command delete deployment")
		client := util.GetClient(ctx)
		namespace := util.GetNamespace(ctx)
		if ctx.NArg() == 0 {
			logrus.Debugf("show help")
			return cli.ShowSubcommandHelp(ctx)
		}
		depl := ctx.Args().First()
		logrus.Debugf("deleting deployment %q from %q", depl, namespace)
		err := client.DeleteDeployment(namespace, depl)
		logrus.WithError(err).Debugf("error while deleting deployment")
		return err
	},
}
