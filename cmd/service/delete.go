package cliserv

import (
	"strings"

	"github.com/containerum/chkit/cmd/util"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

var Delete = &cli.Command{
	Name:        "service",
	Usage:       "call to delete service in specific namespace",
	UsageText:   "chkit delete service service_label [-n namespace]",
	Description: "deletes service in namespace. Aliases: " + strings.Join(aliases, ", "),
	Aliases:     aliases,
	Flags:       util.DeleteFlags,
	Action: func(ctx *cli.Context) error {
		logrus.Debugf("running command delete service")
		client := util.GetClient(ctx)
		namespace := util.GetNamespace(ctx)
		if ctx.NArg() == 0 {
			logrus.Debugf("show help")
			return cli.ShowSubcommandHelp(ctx)
		}
		service := ctx.Args().First()
		logrus.Debugf("deleting service %q from %q", service, namespace)
		err := client.DeleteService(namespace, service)
		logrus.WithError(err).Debugf("error while deleting service")
		return err
	},
}
