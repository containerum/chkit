package cmd

import (
	"github.com/containerum/chkit/cmd/namespace"
	"gopkg.in/urfave/cli.v2"
)

var commandGet = &cli.Command{
	Name: "get",
	Before: func(ctx *cli.Context) error {
		if ctx.Bool("help") {
			return nil
		}
		return setupAll(ctx)
	},
	ArgsUsage: `get ns [namespace_name] --> show namespace data
	         get ns                  --> show namespace list`,
	Action: func(ctx *cli.Context) error {
		return nil
	},
	Subcommands: []*cli.Command{
		namespace.GetNamespace,
	},
}
