package clideployment

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	. "github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/deployment/deplactive"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/animation"
	"github.com/containerum/chkit/pkg/util/trasher"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/urfave/cli.v2"
)

var createDeplConfig = struct {
	File  string
	Force bool
}{}

var Create = &cobra.Command{
	Use:     "deployment",
	Aliases: aliases,
	Run: func(cmd *cobra.Command, args []string) {
		depl := deplactive.DefaultDeployment()
		if cmd.Flag("file").Changed {
			var err error
			depl, err = deplactive.FromFile(createDeplConfig.File)
			if err != nil {
				logrus.WithError(err).Errorf("unable to load deployment data from file %s", createDeplConfig.File)
				fmt.Printf("Unable to load deployment data from file :(\n%v", err)
				os.Exit(1)
			}
		}

		depl, err := deplactive.ConstructDeployment(deplactive.Config{
			Deployment: &depl,
		})
		if err != nil {
			logrus.WithError(err).Errorf("unable to create deployment")
			fmt.Println(err)
			os.Exit(1)
		}
		for {
			_, err := (&activekit.Menu{
				Items: []*activekit.MenuItem{
					{
						Name: "Create deployment",
						Action: func() error {
							if err := Context.Client.CreateDeployment(Context.Namespace, depl); err != nil {
								logrus.WithError(err).Errorf("unable to create deployment %q", depl.Name)
								fmt.Println(err)
							}
							return nil
						},
					},
					{
						Name: "Print to terminal",
						Action: activekit.ActionWithErr(func() error {
							if data, err := depl.RenderYAML(); err != nil {
								return err
							} else {
								fmt.Println(data)
							}
							return nil
						}),
					},
					{
						Name: "Save to file",
						Action: func() error {
							filename, _ := activekit.AskLine("Print filename: ")
							data, err := depl.RenderJSON()
							if err != nil {
								return err
							}
							if err := ioutil.WriteFile(filename, []byte(data), os.ModePerm); err != nil {
								logrus.WithError(err).Errorf("unable to save deployment %q to file", depl.Name)
								fmt.Printf("Unable to save deployment to file :(\n%v", err)
								return nil
							}
							fmt.Printf("OK\n")
							return nil
						},
					},
					{
						Name: "Exit",
						Action: func() error {
							os.Exit(0)
							return nil
						},
					},
				},
			}).Run()
			if err != nil {
				logrus.WithError(err).Errorf("error while menu execution")
				angel.Angel(err)
				os.Exit(1)
			}
		}
	},
}

var CreateOld = &cli.Command{
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
		client := Context.Client
		namespace := Context.Namespace
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
			_, option, _ := activekit.Options("What do you want to do with deployment?", false,
				"Push to server",
				"Print to terminal",
				"Dump to file",
				"Exit")
			switch option {
			case 0:
				anime := &animation.Animation{
					Framerate:      0.3,
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
					fmt.Printf("\n%v\n", err)
				}
			case 1:
				data, _ := depl.RenderYAML()
				w := textWidth(data)
				fmt.Println(strings.Repeat("-", w))
				fmt.Println(data)
				fmt.Println(strings.Repeat("-", w))
			case 2:
				filename, _ := activekit.AskLine("Print filename > ")
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

func init() {
	Create.PersistentFlags().
		StringVar(&createDeplConfig.File, "file", "", "create deployment from file")
	Create.PersistentFlags().
		BoolVarP(&createDeplConfig.Force, "force", "f", false, "create from file without customisation")
}
