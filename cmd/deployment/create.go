package clideployment

import (
	"fmt"

	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/model/deployment/deplactive"
	"github.com/containerum/chkit/pkg/util/activeToolkit"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

var Create = &cli.Command{
	Name:    "deployment",
	Aliases: aliases,
	Action: func(ctx *cli.Context) error {
		client := util.GetClient(ctx)
		namespace := util.GetNamespace(ctx)
		depl, err := deplactive.ConstructDeployment(deplactive.Config{})
		if err != nil {
			logrus.WithError(err).Error("error while creating service")
			fmt.Printf("%v", err)
			return err
		}
		fmt.Println(depl.RenderTable())
		yes, _ := activeToolkit.Yes("Do you want to push deployment to server?")
		if !yes {
			return nil
		}
		err = client.CreateDeployment(namespace, depl)
		if err != nil {
			logrus.WithError(err).Error("unable to create deployment")
			return err
		}
		return nil
	},
}
