package clinamespace

import (
	"strings"
	"time"

	"github.com/containerum/chkit/pkg/model/namespace"
	"github.com/containerum/chkit/pkg/util/animation"
	"github.com/containerum/chkit/pkg/util/trasher"
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
	Description: `shows namespace data or namespace list. Aliases: ` + strings.Join(aliases, ", "),
	Usage:       `shows namespace data or namespace list`,
	UsageText:   "chkit get namespace_name... [-o yaml/json] [-f output_file]",
	Action: func(ctx *cli.Context) error {
		client := util.GetClient(ctx)
		defer util.StoreClient(ctx, client)
		var showItem model.Renderer
		var err error

		anime := &animation.Animation{
			Framerate:      0.2,
			ClearLastFrame: true,
			Source:         trasher.NewSilly(),
		}
		go func() {
			time.Sleep(time.Second)
			anime.Run()
		}()
		switch ctx.NArg() {
		case 1:
			namespaceLabel := ctx.Args().First()
			logrus.Debugf("getting namespace %q", namespaceLabel)
			showItem, err = client.GetNamespace(namespaceLabel)
			if err != nil {
				logrus.Debugf("fatal error: %v", err)
				anime.Stop()
				return err
			}
		default:
			var list namespace.NamespaceList
			logrus.Debugf("getting namespace list")
			list, err := client.GetNamespaceList()
			if err != nil {
				logrus.Debugf("fatal error: %v", err)
				anime.Stop()
				return err
			}
			showItem = list
		}

		anime.Stop()
		err = util.ExportDataCommand(ctx, showItem)
		if err != nil {
			logrus.Debugf("fatal error: %v", err)
		}
		return err
	},
	Flags: util.GetFlags,
}
