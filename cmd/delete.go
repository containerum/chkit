package cmd

import (
	"github.com/containerum/chkit/cmd/deployment"
	"github.com/containerum/chkit/cmd/namespace"
	"github.com/containerum/chkit/cmd/pod"
	"github.com/containerum/chkit/cmd/service"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

var commandDelete = &cli.Command{
	Name:      "delete",
	Usage:     "delete resource",
	UsageText: `chkit delete namespace|deployment|pod|service label`,
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
		clinamespace.Delete,
		clipod.Delete,
		cliserv.Delete,
	},
}
