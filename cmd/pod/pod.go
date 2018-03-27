package clipod

import (
	"fmt"

	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/pod"
	"gopkg.in/urfave/cli.v2"
)

var GetPodAction = &cli.Command{
	Name: "pod",
	Action: func(ctx *cli.Context) error {
		log := util.GetLog(ctx)
		client := util.GetClient(ctx)

		defer func() {
			log := util.GetLog(ctx)
			util.SetClient(ctx, client)
			log.Debugf("writing tokens to disk")
			err := util.SaveTokens(ctx, client.Tokens)
			if err != nil {
				log.Debugf("error while saving tokens: %v", err)
				panic(err)
			}
		}()

		var showItem model.Renderer
		var err error
		switch ctx.NArg() {
		case 0:
			fmt.Println(ctx.Command.UsageText)
			return nil
		case 1:
			namespaceLabel := ctx.Args().First()
			log.Debugf("getting pod list from %q", namespaceLabel)
			showItem, err = client.GetPodList(namespaceLabel)
			if err != nil {
				return err
			}
		default:
			var list pod.PodList
			var gainedPod pod.Pod
			namespaceLabel := ctx.Args().First()
			log.Debugf("getting pods")
			for _, podName := range ctx.Args().Tail() {
				log.Debugf("getting %q", podName)
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
			log.Debugf("fatal error: %v", err)
		}
		return err
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "json",
		},
		&cli.StringFlag{
			Name: "yaml",
		},
	},
}
