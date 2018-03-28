package clipod

import (
	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/pod"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

var GetPodAction = &cli.Command{
	Name:    "pod",
	Aliases: []string{"po", "pods"},
	Action: func(ctx *cli.Context) error {
		client := util.GetClient(ctx)
		defer util.StoreClient(ctx, client)
		var showItem model.Renderer
		var err error

		switch ctx.NArg() {
		case 0:
			namespaceLabel := util.GetNamespace(ctx)
			logrus.Debugf("getting pod list from %q", namespaceLabel)
			showItem, err = client.GetPodList(namespaceLabel)
			if err != nil {
				return err
			}
		case 1:
			namespaceLabel := util.GetNamespace(ctx)
			showItem, err = client.GetPod(namespaceLabel, ctx.Args().First())
			if err != nil {
				return err
			}
		default:
			logrus.Debugf("getting pods")
			gainedList, err := client.GetPodList(util.GetNamespace(ctx))
			var list pod.PodList
			if err != nil {
				return err
			}
			podNames := util.NewSet(ctx.Args().Slice())
			for _, pod := range gainedList {
				if podNames.Have(pod.Name) {
					list = append(list, pod)
				}
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
