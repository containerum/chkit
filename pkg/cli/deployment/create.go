package clideployment

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/model/deployment/deplactive"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Create(ctx *context.Context) *cobra.Command {
	var file string
	var force bool
	var flagCont container.Container
	var flagDepl deployment.Deployment
	var envs []string
	command := &cobra.Command{
		Use:     "deployment",
		Aliases: aliases,
		Short:   "create new deployment",
		Long: `Creates new deployment.
Has an one-line mode, suitable for integration with other tools, and an interactive wizard mode`,
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
					envMap := map[string]string{}
					for _, env := range envs {
						var tokens = strings.SplitN(env, ":", 2)
						if len(tokens) != 2 {
							fmt.Printf("Inavlid env kev-value pair %q", env)
							os.Exit(1)
						}
						var k, v = tokens[0], tokens[1]
						k = strings.TrimSpace(k)
						v = strings.TrimSpace(v)
						v = strings.TrimPrefix(v, "\"")
						v = strings.TrimSuffix(v, "\"")
						envMap[k] = v
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

			var firstItem *activekit.MenuItem
			var created = false
			if activekit.YesNo("Are you sure?") {
				if err := ctx.Client.CreateDeployment(ctx.Namespace, depl); err != nil {
					logrus.WithError(err).Errorf("unable to create deployment %q", depl.Name)
					activekit.Attention(err.Error())
					os.Exit(1)
				}
				fmt.Printf("Congratulations! Deployment %q created!\n", depl.Name)
				created = true
				firstItem = &activekit.MenuItem{
					Label: "Push changes to server",
					Action: func() error {
						if activekit.YesNo("Are you sure?") {
							err := ctx.Client.ReplaceDeployment(ctx.Namespace, depl)
							if err != nil {
								logrus.WithError(err).Errorf("unable to update deployment %q", depl.Name)
								fmt.Println(err)
								return nil
							}
							fmt.Printf("Congratulations! Deployment %q updated!\n", depl.Name)
						}
						return nil
					},
				}
			} else {
				firstItem = &activekit.MenuItem{
					Label: "Create deployment on server",
					Action: func() error {
						if activekit.YesNo("Are you sure?") {
							err := ctx.Client.CreateDeployment(ctx.Namespace, depl)
							if err != nil {
								logrus.WithError(err).Errorf("unable to update deployment %q", depl.Name)
								fmt.Println(err)
								return nil
							}
							fmt.Printf("Congratulations! Deployment %q created!\n", depl.Name)
							created = true
						}
						return nil
					},
				}
			}
			for {
				deploymentMenu(ctx, depl, firstItem).Run()
				if created {
					firstItem = &activekit.MenuItem{
						Label: "Push changes to server",
						Action: func() error {
							if activekit.YesNo("Are you sure?") {
								err := ctx.Client.ReplaceDeployment(ctx.Namespace, depl)
								if err != nil {
									logrus.WithError(err).Errorf("unable to update deployment %q", depl.Name)
									fmt.Println(err)
									return nil
								}
								fmt.Printf("Congratulations! Deployment %q updated!\n", depl.Name)
							}
							return nil
						},
					}
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
		StringSliceVar(&envs, "env", []string{}, "container env variable in KEY0:VALUE0 KEY1:VALUE1 format")
	return command
}

func deploymentMenu(ctx *context.Context, depl deployment.Deployment, firstItem *activekit.MenuItem) *activekit.Menu {
	return &activekit.Menu{
		Items: []*activekit.MenuItem{
			firstItem,
			{
				Label: "Edit deployment",
				Action: func() error {
					var err error
					depl, err = deplactive.ReplaceWizard(deplactive.Config{
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
	}
}
