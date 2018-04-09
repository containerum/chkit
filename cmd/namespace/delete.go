package clinamespace

import (
	"strings"

	"github.com/containerum/chkit/cmd/cmdutil"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

var Delete = &cli.Command{
	Name:        "namespace",
	Usage:       "call to delete namespace",
	Description: "delete namespace deletes namespace with name, provided by first arg. Aliases: " + strings.Join(aliases, ", "),
	UsageText:   "chkit delete namespace",
	Aliases:     aliases,
	Flags:       cmdutil.DeleteFlags,
	Action: func(ctx *cli.Context) error {
		logrus.Debugf("running command delete namespace")
		client := cmdutil.GetClient(ctx)
		if ctx.NArg() == 0 {
			logrus.Debugf("show help")
			return cli.ShowSubcommandHelp(ctx)
		}
		namespace := ctx.Args().First()
		logrus.Debugf("deleting namespace %q", namespace)
		err := client.DeleteNamespace(namespace)
		return err
	},
}
