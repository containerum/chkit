package clideployment

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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
		if yes {
			err = client.CreateDeployment(namespace, depl)
			if err != nil {
				logrus.WithError(err).Error("unable to create deployment")
				fmt.Println(err)
			}
		}
		_, option, _ := activeToolkit.Options("Do you want to dump deploymnent?", false,
			"Yes, to file",
			"Yes, to stdout",
			"No")

		switch option {
		case 0:
			filename, _ := activeToolkit.AskLine("Print filename > ")
			if strings.TrimSpace(filename) == "" {
				return nil
			}
			depl.ToKube()
			data, _ := depl.MarshalJSON()
			err := ioutil.WriteFile(filename, data, os.ModePerm)
			if err != nil {
				logrus.WithError(err).Error("unable to write deployment to file")
				fmt.Println(err)
				return nil
			}
		case 1:
			data, _ := depl.RenderYAML()
			fmt.Println(data)
			return nil
		default:
			return nil
		}
		return nil
	},
}
