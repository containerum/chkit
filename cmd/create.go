package cmd

import (
	"github.com/containerum/chkit/cmd/service"
	"gopkg.in/urfave/cli.v2"
)

var CommandCreate = &cli.Command{
	Name: "create",
	Action: func(ctx *cli.Context) error {
		return cli.ShowSubcommandHelp(ctx)
	},
	Subcommands: []*cli.Command{
		cliserv.Create,
	},
}
