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
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Replace(ctx *context.Context) *cobra.Command {
	var file string
	var force bool
	var flagCont container.Container
	var flagDepl deployment.Deployment
	var envs []string
	command := &cobra.Command{
		Use:     "deployment",
		Aliases: aliases,
		Short:   "replace deployment",
		Long: `Replaces deployment.
Has an one-line mode, suitable for integration with other tools, and an interactive wizard mode`,
		Run: func(cmd *cobra.Command, args []string) {
			depl := deplactive.DefaultDeployment()
			switch len(args) {
			case 1:
				depl.Name = args[0]
			default:
				cmd.Help()
			}
			if cmd.Flag("file").Changed {
				var err error
				depl, err = deplactive.FromFile(file)
				if err != nil {
					logrus.WithError(err).Errorf("unable to load deployment data from file %s", file)
					fmt.Printf("Unable to load deployment data from file :(\n%v", err)
					os.Exit(1)
				}
			} else {
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

				oldDepl, err := ctx.Client.GetDeployment(ctx.Namespace, depl.Name)
				if err != nil {
					activekit.Attention("unable to get deploment %q: %v", depl.Name, err)
					os.Exit(1)
				}
				if !cmd.Flag("replicas").Changed {
					flagDepl.Replicas = oldDepl.Replicas
				}
				if !cmd.Flag("image").Changed {
					flagDepl.Containers = oldDepl.Containers
				} else {
					flagDepl.Containers = []container.Container{flagCont}
				}

				depl = flagDepl
			}
			if cmd.Flag("force").Changed {
				if len(args) != 1 {
					cmd.Help()
					return
				}
				depl.Name = args[0]
				if err := deplactive.ValidateDeployment(depl); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(depl.RenderTable())
				if err := ctx.Client.ReplaceDeployment(ctx.Namespace, depl); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println("OK")
				return
			} else {
				if len(args) == 0 {
					list, err := ctx.Client.GetDeploymentList(ctx.Namespace)
					if err != nil {
						activekit.Attention(err.Error())
						os.Exit(1)
					}
					var menu []*activekit.MenuItem
					for _, d := range list {
						menu = append(menu, &activekit.MenuItem{
							Label: d.Name,
							Action: func(d deployment.Deployment) func() error {
								return func() error {
									depl = d
									return nil
								}
							}(d),
						})
					}
					(&activekit.Menu{
						Title: "Choose deployment to replace",
						Items: menu,
					}).Run()
				} else {
					var err error
					depl, err = ctx.Client.GetDeployment(ctx.Namespace, args[0])
					if err != nil {
						activekit.Attention(err.Error())
						os.Exit(1)
					}
				}
			}
			depl, err := deplactive.ReplaceWizard(deplactive.Config{
				Deployment: &depl,
			})
			if err != nil {
				logrus.WithError(err).Errorf("unable to replace deployment")
				fmt.Println(err)
				os.Exit(1)
			}
			for {
				_, err := (&activekit.Menu{
					Items: []*activekit.MenuItem{
						{
							Label: fmt.Sprintf("Update deployment %q on server", depl.Name),
							Action: func() error {
								fmt.Println(depl.RenderTable())
								if activekit.YesNo(fmt.Sprintf("Are you sure you want to update deployment %q on server?", depl.Name)) {
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
						},
						{
							Label: "Edit deployment",
							Action: func() error {
								var err error
								depl, err = deplactive.ReplaceWizard(deplactive.Config{
									Deployment: &depl,
								})
								if err != nil {
									logrus.WithError(err).Errorf("unable to update deployment")
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
								if filename == "" {
									return nil
								}
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
		IntVar(&flagDepl.Replicas, "replicas", 1, "replicas, optional")
	command.PersistentFlags().
		StringVar(&flagCont.Name, "container-name", "", "container name, equal to image name by default")
	command.PersistentFlags().
		StringVar(&flagCont.Image, "image", "", "container image, optional")
	command.PersistentFlags().
		UintVar(&flagCont.Limits.Memory, "memory", 256, "container memory limit im Mb, optional")
	command.PersistentFlags().
		UintVar(&flagCont.Limits.CPU, "cpu", 200, "container CPU limit in mCPU, optional")
	command.PersistentFlags().
		StringArrayVar(&envs, "env", nil, "container env variable in KEY0:VALUE0 KEY1:VALUE1 format")
	return command
}
