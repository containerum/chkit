package clinamespace

import (
	"fmt"
	"strings"
	"time"

	"github.com/containerum/chkit/cmd/cmdutil"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/namespace"
	"github.com/containerum/chkit/pkg/util/animation"
	"github.com/containerum/chkit/pkg/util/trasher"
	"github.com/sirupsen/logrus"
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
		client := cmdutil.GetClient(ctx)
		defer cmdutil.StoreClient(ctx, client)
		var showItem model.Renderer
		var err error

		anime := &animation.Animation{
			Framerate:      0.5,
			ClearLastFrame: true,
			Source:         trasher.NewSilly(),
		}
		go func() {
			time.Sleep(4 * time.Second)
			anime.Run()
		}()

		err = func() error {
			defer anime.Stop()
			switch ctx.NArg() {
			case 1:
				namespaceLabel := ctx.Args().First()
				logrus.Debugf("getting namespace %q", namespaceLabel)
				showItem, err = client.GetNamespace(namespaceLabel)
				if err != nil {
					logrus.WithError(err).Errorf("unable to get namespace %q", namespaceLabel)
					fmt.Printf("Error hile getting namespace %q: %v\n", namespaceLabel, err)
					return err
				}
			default:
				var list namespace.NamespaceList
				logrus.Debugf("getting namespace list")
				list, err := client.GetNamespaceList()
				if err != nil {
					logrus.WithError(err).Errorf("unable to get namespace list")
					return err
				}
				defaultNamespace := cmdutil.GetNamespace(ctx)
				fmt.Printf("Using %q as default namespace\n", defaultNamespace)
				showItem = list
			}
			return nil
		}()
		if err != nil {
			return err
		}
		logrus.Debugf("List recieved")
		err = cmdutil.ExportDataCommand(ctx, showItem)
		if err != nil {
			logrus.Debugf("fatal error: %v", err)
		}
		return err
	},
	Flags: cmdutil.GetFlags,
}
