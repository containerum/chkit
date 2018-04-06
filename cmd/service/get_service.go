package cliserv

import (
	"strings"
	"time"

	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/animation"
	"github.com/containerum/chkit/pkg/util/trasher"
	"gopkg.in/urfave/cli.v2"
)

var aliases = []string{"srv", "services", "svc"}

var GetService = &cli.Command{
	Name:        "service",
	Usage:       "shows service info",
	UsageText:   "chkit get service service_label [-o yaml/json] [-f output_file]",
	Description: "shows service info. Aliases: " + strings.Join(aliases, ", "),
	Aliases:     aliases,
	Action: func(ctx *cli.Context) error {
		client := util.GetClient(ctx)
		defer util.StoreClient(ctx, client)
		var show model.Renderer
		var err error

		anime := &animation.Animation{
			Framerate:      0.5,
			ClearLastFrame: true,
			Source:         trasher.NewSilly(),
		}
		go func() {
			time.Sleep(time.Second)
			anime.Run()
		}()

		switch ctx.NArg() {
		case 0:
			namespace := util.GetNamespace(ctx)
			list, err := client.GetServiceList(namespace)
			if err != nil {
				anime.Stop()
				return err
			}
			show = list
		case 1:
			namespace := util.GetNamespace(ctx)
			show, err = client.GetService(namespace, ctx.Args().First())
			if err != nil {
				anime.Stop()
				return err
			}
		default:
			namespace := util.GetNamespace(ctx)
			servicesNames := util.NewSet(ctx.Args().Slice())
			gainedList, err := client.GetServiceList(namespace)
			if err != nil {
				anime.Stop()
				return err
			}
			var list service.ServiceList
			for _, serv := range gainedList {
				if servicesNames.Have(serv.Name) {
					list = append(list, serv)
				}
			}
			show = list
		}
		anime.Stop()
		return util.ExportDataCommand(ctx, show)
	},
	Flags: util.GetFlags,
}
