package clinamespace

import (
	"fmt"

	"github.com/containerum/chkit/pkg/model/namespace"

	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/urfave/cli.v2"
)

// GetNamespace -- commmand 'get' entity data
var GetNamespace = &cli.Command{
	Name:        "ns",
	Description: `show namespace or namespace list`,
	Usage:       `Shows namespace data or namespace list`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "json",
		},
	},
	Action: func(ctx *cli.Context) error {
		log := util.GetLog(ctx)
		client := util.GetClient(ctx)
		if ctx.NArg() == 0 {

		}

		var showItem model.Renderer
		var err error
		switch ctx.NArg() {
		case 1:
			namespaceLabel := ctx.Args().First()
			log.Debugf("getting namespace %q", namespaceLabel)
			showItem, err = client.GetNamespace(namespaceLabel)
			if err != nil {
				log.Debugf("fatal error: %v", err)
				return err
			}
		default:
			var list namespace.NamespaceList
			log.Debugf("getting namespace list")
			list, err := client.GetNamespaceList()
			if err != nil {
				log.Debugf("fatal error: %v", err)
				return err
			}
			showItem = list
		}
		switch {
		case ctx.IsSet("json"):
		case ctx.IsSet("yaml"):
		default:
			fmt.Println(showItem.RenderTable())
		}
		return nil
	},
}
