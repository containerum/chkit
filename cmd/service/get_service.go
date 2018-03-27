package cliserv

import (
	"fmt"

	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/service"
	"gopkg.in/urfave/cli.v2"
)

var GetService = &cli.Command{
	Name: "service",
	Action: func(ctx *cli.Context) error {
		client := util.GetClient(ctx)

		defer func() {
			log := util.GetLog(ctx)
			util.SetClient(ctx, client)
			log.Debugf("writing tokens to disk %s", client.Tokens)
			err := util.SaveTokens(ctx, client.Tokens)
			if err != nil {
				log.Debugf("error while saving tokens: %v", err)
				panic(err)
			}
		}()

		var show model.Renderer
		switch ctx.NArg() {
		case 0:
			fmt.Println(ctx.Command.Usage)
			return nil
		case 1:
			namespace := ctx.Args().First()
			list, err := client.GetServiceList(namespace)
			if err != nil {
				return err
			}
			show = list
		default:
			namespace := ctx.Args().First()
			servicesNames := ctx.Args().Tail()
			var list service.ServiceList
			for _, servName := range servicesNames {
				serv, err := client.GetService(namespace, servName)
				if err != nil {
					return err
				}
				list = append(list, serv)
			}
			show = list
		}
		return util.WriteData(ctx, show)
	},
	After: func(ctx *cli.Context) error {
		client := util.GetClient(ctx)
		log := util.GetLog(ctx)
		log.Debugf("writing tokens to disk")
		return util.SaveTokens(ctx, client.Tokens)
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
