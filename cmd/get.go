package cmd

import (
	"github.com/containerum/chkit/cmd/deployment"
	"github.com/containerum/chkit/cmd/namespace"
	"github.com/containerum/chkit/cmd/pod"
	"github.com/containerum/chkit/cmd/service"
	"gopkg.in/urfave/cli.v2"
)

var commandGet = &cli.Command{
	Name: "get",
	Before: func(ctx *cli.Context) error {
		if ctx.Bool("help") {
			return cli.ShowSubcommandHelp(ctx)
		}
		return setupAll(ctx)
	},
	ArgsUsage: `get ns [namespace_name] --> show namespace data
	         get ns                  --> show namespace list`,
	Action: func(ctx *cli.Context) error {
		cli.ShowCommandHelpAndExit(ctx, "get", 0)
		return nil
	},
	Subcommands: []*cli.Command{
		clinamespace.GetNamespace,
		clipod.GetPodAction,
		cliserv.GetService,
		clideployment.GetDeployment,
	},
}
