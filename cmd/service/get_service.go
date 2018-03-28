package cliserv

import (
	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/service"
	"gopkg.in/urfave/cli.v2"
)

var GetService = &cli.Command{
	Name:    "service",
	Aliases: []string{"srv", "services", "svc"},
	Action: func(ctx *cli.Context) error {
		client := util.GetClient(ctx)
		defer util.StoreClient(ctx, client)
		var show model.Renderer
		var err error
		switch ctx.NArg() {
		case 0:
			namespace := util.GetNamespace(ctx)
			list, err := client.GetServiceList(namespace)
			if err != nil {
				return err
			}
			show = list
		case 1:
			namespace := util.GetNamespace(ctx)
			show, err = client.GetService(namespace, ctx.Args().First())
			if err != nil {
				return err
			}
		default:
			namespace := util.GetNamespace(ctx)
			servicesNames := util.NewSet(ctx.Args().Slice())
			gainedList, err := client.GetServiceList(namespace)
			if err != nil {
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
		return util.WriteData(ctx, show)
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
