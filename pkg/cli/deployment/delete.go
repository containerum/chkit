package clideployment

import (
	"strings"

	"github.com/containerum/chkit/pkg/context"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

var DeleteDeployment = &cli.Command{
	Name:        "deployment",
	Usage:       "call to delete deployment in specific namespace",
	UsageText:   "chkit delete deployment deployment_label [-n namespace]",
	Description: "deletes deployment. Aliases: " + strings.Join(aliases, ", "),
	Aliases:     aliases,
	Action: func(ctx *cli.Context) error {
		logrus.Debugf("running command delete deployment")
		client := context.GlobalContext.Client
		if ctx.NArg() == 0 {
			logrus.Debugf("show help")
			return cli.ShowSubcommandHelp(ctx)
		}
		depl := ctx.Args().First()
		logrus.Debugf("deleting deployment %q from %q", depl, context.GlobalContext.Namespace)
		err := client.DeleteDeployment(context.GlobalContext.Namespace, depl)
		logrus.WithError(err).Debugf("error while deleting deployment")
		return err
	},
}
