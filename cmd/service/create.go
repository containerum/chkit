package cliserv

import (
	"fmt"

	"github.com/containerum/chkit/pkg/model/service/servactive"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

var Create = &cli.Command{
	Name:    "service",
	Aliases: aliases,
	Action: func(ctx *cli.Context) error {
		serv, err := servactive.RunInteractveConstructor(servactive.ConstructorConfig{
			Force: ctx.Bool("force"),
		})
		if err != nil {
			logrus.WithError(err).Debugf("error while constructing service")
			return err
		}
		fmt.Println(serv.RenderTable())
		return nil
	},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Value:   false,
		},
	},
}
