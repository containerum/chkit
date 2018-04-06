package clideployment

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/model/deployment/deplactive"
	"github.com/containerum/chkit/pkg/util/activeToolkit"
	"github.com/containerum/chkit/pkg/util/animation"
	"github.com/containerum/chkit/pkg/util/trasher"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

var Create = &cli.Command{
	Name:    "deployment",
	Aliases: aliases,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "file with deployment data",
		},
	},
	Action: func(ctx *cli.Context) error {
		client := util.GetClient(ctx)
		namespace := util.GetNamespace(ctx)
		deplConfig := deplactive.Config{}
		if ctx.IsSet("file") {
			deploymentFile := ctx.String("file")
			depl, err := deplactive.FromFile(deploymentFile)
			if err != nil {
				logrus.WithError(err).
					Errorf("unable to read deployment data from %q", deploymentFile)
				fmt.Printf("Unable to read data from %q: %v\n", deploymentFile, err)
				return err
			}
			deplConfig.Deployment = &depl
		}
		depl, err := deplactive.ConstructDeployment(deplConfig)
		if err != nil {
			logrus.WithError(err).Error("error while creating deployment")
			fmt.Printf("%v\n", err)
			return err
		}
		fmt.Println(depl.RenderTable())
		for {
			_, option, _ := activeToolkit.Options("What do you want to do with deployment?", false,
				"Push to server",
				"Print to terminal",
				"Dump to file",
				"Exit")
			switch option {
			case 0:
				anime := &animation.Animation{
					Framerate:      0.5,
					ClearLastFrame: true,
					Source:         trasher.NewSilly(),
				}
				go func() {
					time.Sleep(time.Second)
					anime.Run()
				}()
				go anime.Run()
				err = client.CreateDeployment(namespace, depl)
				anime.Stop()
				if err != nil {
					logrus.WithError(err).Error("unable to create deployment")
					fmt.Println(err)
				}
			case 1:
				data, _ := depl.RenderYAML()
				w := textWidth(data)
				fmt.Println(strings.Repeat("-", w))
				fmt.Println(data)
				fmt.Println(strings.Repeat("-", w))
			case 2:
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
				}
			default:
				return nil
			}
		}
	},
}

func textWidth(text string) int {
	width := 0
	for _, line := range strings.Split(text, "\n") {
		if len(line) > width {
			width = len(line)
		}
	}
	return width
}
