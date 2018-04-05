package clideployment

import (
	"fmt"

	"github.com/containerum/chkit/pkg/model/deployment/deplactive"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

var Create = &cli.Command{
	Name:    "deployment",
	Aliases: aliases,
	Action: func(ctx *cli.Context) error {
		depl, err := deplactive.ConstructDeployment(deplactive.Config{})
		if err != nil {
			logrus.WithError(err).Error("error while creating service")
			fmt.Printf("%v", err)
			return err
		}
		fmt.Println(depl.RenderTable())
		return nil
	},
}
