package namespace

import (
	"fmt"

	"github.com/containerum/chkit/cmd/util"
	"gopkg.in/urfave/cli.v2"
)

var GetNamespace = &cli.Command{
	Name:        "ns",
	Description: `show namespace or namespace list`,
	Usage:       `chkit get ns newton\nchkit get ns`,
	Action: func(ctx *cli.Context) error {
		client := util.GetClient(ctx)
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
}
