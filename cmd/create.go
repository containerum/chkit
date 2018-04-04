package cmd

import (
	"github.com/containerum/chkit/cmd/service"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

var CommandCreate = &cli.Command{
	Name: "create",
	Before: func(ctx *cli.Context) error {
		logrus.Debugf("start create action")
		return setupAll(ctx)
	},
	Action: func(ctx *cli.Context) error {
		return cli.ShowSubcommandHelp(ctx)
	},
	Subcommands: []*cli.Command{
		cliserv.Create,
	},
}
