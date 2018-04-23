package clideployment

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/model/deployment/deplactive"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/containerum/chkit/pkg/util/pairs"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Create(ctx *context.Context) *cobra.Command {
	var file string
	var force bool
	var flagCont container.Container
	var flagDepl deployment.Deployment
	var envs string
	command := &cobra.Command{
		Use:     "deployment",
		Aliases: aliases,
		Short:   "create new deployment",
		Long: `Creates new deployment.
Has an one-line, suitable for integration with other tools, and an interactive wizard mode`,
		Run: func(cmd *cobra.Command, args []string) {
			depl := deplactive.DefaultDeployment()
			if cmd.Flag("file").Changed {
				var err error
				depl, err = deplactive.FromFile(file)
				if err != nil {
					logrus.WithError(err).Errorf("unable to load deployment data from file %s", file)
					fmt.Printf("Unable to load deployment data from file :(\n%v", err)
					os.Exit(1)
				}
			} else if cmd.Flag("force").Changed {
				if cmd.Flag("env").Changed {
					envMap, err := pairs.ParseMap(envs, ":")
					if err != nil {
						fmt.Printf("invalid env flag\n")
						os.Exit(1)
					}
					for k, v := range envMap {
						flagCont.Env = append(flagCont.Env, model.Env{
							Name:  k,
							Value: v,
						})
					}
				}
				if flagCont.Name == "" {
					flagCont.Name = namegen.Aster() + "-" + flagCont.Image
				}
				flagDepl.Containers = []container.Container{flagCont}
				depl = flagDepl
			}
			if cmd.Flag("force").Changed {
				if err := deplactive.ValidateDeployment(depl); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(depl.RenderTable())
				if err := ctx.Client.CreateDeployment(ctx.Namespace, depl); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println("OK")
				return
			}
			depl, err := deplactive.Wizard(deplactive.Config{
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
							Label: "Push deployment to server",
							Action: func() error {
								err := ctx.Client.CreateDeployment(ctx.Namespace, depl)
								if err != nil {
									logrus.WithError(err).Errorf("unable to create deployment %q", depl.Name)
									fmt.Println(err)
									return nil
								}
								fmt.Printf("Congratulations! Deployment %q created!\n", depl.Name)
								return nil
							},
						},
						{
							Label: "Edit deployment",
							Action: func() error {
								var err error
								depl, err = deplactive.Wizard(deplactive.Config{
									Deployment: &depl,
								})
								if err != nil {
									logrus.WithError(err).Errorf("unable to create deployment")
									fmt.Println(err)
									os.Exit(1)
								}
								return nil
							},
						},
						{
							Label: "Print to terminal",
							Action: activekit.ActionWithErr(func() error {
								if data, err := depl.RenderYAML(); err != nil {
									return err
								} else {
									upBorders := strings.Repeat("_", text.Width(data))
									downBorders := strings.Repeat("_", text.Width(data))
									fmt.Printf("%s\n\n%s\n%s\n", upBorders, data, downBorders)
								}
								return nil
							}),
						},
						{
							Label: "Save to file",
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
							Label: "Exit",
							Action: func() error {
								if yes, _ := activekit.Yes("Are you sure you want to exit?"); yes {
									os.Exit(0)
								}
								return nil
							},
						},
					},
				}).Run()
				if err != nil {
					logrus.WithError(err).Errorf("error while menu execution")
					angel.Angel(ctx, err)
					os.Exit(1)
				}
			}
		},
	}
	command.PersistentFlags().
		StringVar(&file, "file", "", "create deployment from file")
	command.PersistentFlags().
		BoolVarP(&force, "force", "f", false, "suppress confirmation")

	command.PersistentFlags().
		StringVar(&flagDepl.Name, "deployment-name", namegen.Color()+"-"+namegen.Aster(), "deployment name, optional")
	command.PersistentFlags().
		IntVar(&flagDepl.Replicas, "replicas", 1, "replicas, optional")
	command.PersistentFlags().
		StringVar(&flagCont.Name, "container-name", "", "container name, equal to image name by default")
	command.PersistentFlags().
		StringVar(&flagCont.Image, "image", "", "container image, required")
	command.PersistentFlags().
		UintVar(&flagCont.Limits.Memory, "memory", 256, "container memory limit im Mb, optional")
	command.PersistentFlags().
		UintVar(&flagCont.Limits.CPU, "cpu", 200, "container CPU limit in mCPU, optional")
	command.PersistentFlags().
		StringSliceVar(&flagCont.Commands, "commands", nil, "container commands")
	command.PersistentFlags().
		StringVar(&envs, "env", "", "container env variable in KEY0:VALUE0 KEY1:VALUE1 format")
	return command
}
