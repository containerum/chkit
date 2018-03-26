package cmd

import (
	"github.com/containerum/chkit/cmd/deployment"
	"github.com/containerum/chkit/cmd/namespace"
	"github.com/containerum/chkit/cmd/pod"
	"github.com/containerum/chkit/cmd/service"
	"github.com/containerum/chkit/cmd/util"
	"gopkg.in/urfave/cli.v2"
)

var commandGet = &cli.Command{
	Name: "get",
	Before: func(ctx *cli.Context) error {
		if ctx.Bool("help") {
			return cli.ShowSubcommandHelp(ctx)
		}
		if err := setupLog(ctx); err != nil {
			return err
		}
		return setupAll(ctx)
	},
	ArgsUsage: `get ns [namespace_name] --> show namespace data
	         get ns                  --> show namespace list`,
	Action: func(ctx *cli.Context) error {
		return ctx.App.Command("help").Run(ctx)
	},
	After: func(ctx *cli.Context) error {
		client := util.GetClient(ctx)
		return util.SaveTokens(ctx, client.Tokens)
	},
	Subcommands: []*cli.Command{
		clinamespace.GetNamespace,
		clipod.GetPodAction,
		cliserv.GetService,
		clideployment.GetDeployment,
	},
}
