package cmd

import (
	"github.com/containerum/chkit/cmd/namespace"
	"gopkg.in/urfave/cli.v2"
)

var commandGet = &cli.Command{
	Name:   "get",
	Before: setupAll,
	Action: func(ctx *cli.Context) error {
		return nil
	},
	Subcommands: []*cli.Command{
		namespace.GetNamespace,
	},
}
