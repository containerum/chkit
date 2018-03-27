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
		default:
			var list pod.PodList
			var gainedPod pod.Pod
			namespaceLabel := util.GetNamespace(ctx)
			logrus.Debugf("getting pods")
			for _, podName := range ctx.Args().Slice() {
				logrus.Debugf("getting %q", podName)
				gainedPod, err = client.GetPod(namespaceLabel, podName)
				if err != nil {
					return err
				}
				list = append(list, gainedPod)
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
			Aliases: []string{"f"},
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
		},
	},
}
