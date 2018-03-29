package cmd

import (
	"github.com/containerum/chkit/cmd/deployment"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

var commandDelete = &cli.Command{
	Name: "delete",
	Before: func(ctx *cli.Context) error {
		logrus.Debugf("start delete action")
		return setupAll(ctx)
	},
	Action: func(ctx *cli.Context) error {
		logrus.Debugf("delete main action")
		return cli.ShowSubcommandHelp(ctx)
	},
	Subcommands: []*cli.Command{
		clideployment.DeleteDeployment,
	},
}
