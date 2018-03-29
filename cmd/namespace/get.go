package clinamespace

import (
	"github.com/containerum/chkit/pkg/model/namespace"
	"github.com/sirupsen/logrus"

	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/urfave/cli.v2"
)

var aliases = []string{"ns", "namespaces"}

// GetNamespace -- commmand 'get' entity data
var GetNamespace = &cli.Command{
	Name:        "namespace",
	Aliases:     aliases,
	Description: `shows namespace data or namespace list. Aliases: ns, namespaces`,
	Usage:       `shows namespace data or namespace list`,
	Action: func(ctx *cli.Context) error {
		client := util.GetClient(ctx)
		defer util.StoreClient(ctx, client)
		var showItem model.Renderer
		var err error

		switch ctx.NArg() {
		case 1:
			namespaceLabel := ctx.Args().First()
			logrus.Debugf("getting namespace %q", namespaceLabel)
			showItem, err = client.GetNamespace(namespaceLabel)
			if err != nil {
				logrus.Debugf("fatal error: %v", err)
				return err
			}
		default:
			var list namespace.NamespaceList
			logrus.Debugf("getting namespace list")
			list, err := client.GetNamespaceList()
			if err != nil {
				logrus.Debugf("fatal error: %v", err)
				return err
			}
			showItem = list
		}
		err = util.WriteData(ctx, showItem)
		if err != nil {
			logrus.Debugf("fatal error: %v", err)
		}
		return err
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "file",
			Usage:   "file to write output",
			Aliases: []string{"f"},
		},
		&cli.StringFlag{
			Name:    "output",
			Usage:   "define output formats: yaml, json",
			Aliases: []string{"o"},
		},
	},
}
