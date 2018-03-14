package cmd

import (
	"fmt"

	"gopkg.in/urfave/cli.v2"
)

var commandGet = &cli.Command{
	Name: "get",
	Action: func(ctx *cli.Context) error {
		return nil
	},
	Subcommands: []*cli.Command{
		&cli.Command{
			Name:        "ns",
			Description: `show namespace or namespace list`,
			Usage:       `chkit get ns newton\nchkit get ns`,
			Action: func(ctx *cli.Context) error {
				// INIT
				if err := setupAll(ctx); err != nil {
					return err
				}
				// END INIT
				client := getClient(ctx)
				if ctx.NArg() > 0 {
					name := ctx.Args().First()
					ns, err := client.GetNamespace(name)
					if err != nil {
						return err
					}
					fmt.Println(ns.RenderTable())
				}
				return nil
			},
		},
	},
}
