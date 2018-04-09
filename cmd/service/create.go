package cliserv

import (
	"fmt"

	"github.com/containerum/chkit/cmd/cmdutil"
	"github.com/containerum/chkit/pkg/model/service/servactive"
	"github.com/containerum/chkit/pkg/util/activeToolkit"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

var Create = &cli.Command{
	Name:    "service",
	Aliases: aliases,
	Action: func(ctx *cli.Context) error {
		client := cmdutil.GetClient(ctx)
		ns := cmdutil.GetNamespace(ctx)
		deplList, err := client.GetDeploymentList(ns)
		if err != nil {
			logrus.WithError(err).Errorf("error while gettin deployment list")
			return err
		}
		if len(deplList) == 0 {
			fmt.Printf("You have no deployments to create service :(")
			return nil
		}
		deployments := make([]string, 0, len(deplList))
		for _, depl := range deplList {
			deployments = append(deployments, depl.Name)
		}
		list, err := servactive.RunInteractveConstructor(servactive.ConstructorConfig{
			Force:       ctx.Bool("force"),
			Deployments: deployments,
		})
		switch err {
		case servactive.ErrUserStoppedSession, nil:
			// pass
		default:
			logrus.WithError(err).Debugf("error while constructing service")
			return err
		}
		if len(list) == 0 {
			fmt.Printf("You didn't create any services\n")
			return nil
		}
		fmt.Println(list.RenderTable())
		if yes, _ := activeToolkit.Yes("Do you want to push services to server?"); yes {
			for _, serv := range list {
				if err := client.CreateService(ns, serv); err != nil {
					logrus.WithError(err).Errorf("unable to create service %q", serv.Name)
					return err
				}
			}
		}
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
