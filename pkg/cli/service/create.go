package cliserv

import (
	"strings"

	"fmt"
	"os"

	"io/ioutil"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/service/servactive"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Create(ctx *context.Context) *cobra.Command {
	var createServiceConfig = struct {
		File  string
		Force bool
	}{
		File: "-",
	}
	command := &cobra.Command{
		Use:     "service",
		Aliases: aliases,
		Short:   "create service",
		Long:    "create service for provided pod in provided namespace. Aliases: " + strings.Join(aliases, ", "),
		Run: func(cmd *cobra.Command, args []string) {
			logrus.WithField("command", "create service").Debugf("start service creation")
			depList, err := context.GlobalContext.Client.GetDeploymentList(context.GlobalContext.Namespace)
			if err != nil {
				logrus.WithError(err).Errorf("unable to get deployment list")
				fmt.Println("Unable to get deployment list :(")
			}
			wizardConfig := servactive.ConstructorConfig{
				Deployments: depList.Names(),
			}
			if cmd.Flag("file").Changed && createServiceConfig.File != "" && createServiceConfig.File != "-" {
				logrus.Debugf("loading service from file")
				serv, err := servactive.FromFile(createServiceConfig.File)
				if err != nil {
					logrus.WithError(err).Errorf("unable to load service from file")
					activekit.Attention(err.Error())
					os.Exit(1)
				}
				if createServiceConfig.Force {
					if err := ctx.Client.CreateService(ctx.Namespace, serv); err != nil {
						logrus.WithError(err).Errorf("unable to create service %q in namespace %q", serv.Name, ctx.Namespace)
						activekit.Attention(err.Error())
						os.Exit(1)
					}
					fmt.Printf("Service %q created\n", serv.Name)
					return
				}
				wizardConfig.Service = &serv
			}
			service, err := servactive.Wizard(wizardConfig)
			if err != nil {
				logrus.WithError(err).Errorf("unable to create service")
				fmt.Println("Unable to create service :(")
				os.Exit(1)
			}
			fmt.Println(service.RenderTable())
			for {
				(&activekit.Menu{
					Items: []*activekit.MenuItem{
						{
							Label: "Push service to server",
							Action: func() error {
								if activekit.YesNo("Are you sure? [Y/N]: ") {
									if err := ctx.Client.CreateService(ctx.Namespace, service); err != nil {
										logrus.WithError(err).Errorf("unable to create service %q in namespace %q", service.Name, ctx.Namespace)
										activekit.Attention(err.Error())
										return nil
									}
								}
								logrus.WithError(err).Errorf("service %q in namespace %q created", service.Name, ctx.Namespace)
								fmt.Printf("Service %q created\n", service.Name)
								return nil
							},
						},
						{
							Label: "Edit service",
							Action: func() error {
								s, err := servactive.Wizard(servactive.ConstructorConfig{
									Service:     &service,
									Deployments: depList.Names(),
								})
								if err != nil {
									logrus.WithError(err).Errorf("error while interactive service creation")
									activekit.Attention(err.Error())
									os.Exit(1)
								}
								service = s
								return nil
							},
						},
						{
							Label: "Print to terminal",
							Action: func() error {
								data, err := service.RenderYAML()
								if err != nil {
									logrus.WithError(err).Errorf("unable to render service to yaml")
									activekit.Attention(err.Error())
								}
								border := strings.Repeat("_", text.Width(data))
								fmt.Printf("%s\n%s\n%s\n", border, data, border)
								return nil
							},
						},
						{
							Label: "Save to file",
							Action: func() error {
								logrus.Debugf("saving service to file")
								data, err := service.RenderJSON()
								if err != nil {
									logrus.WithError(err).Errorf("unable to render service to json")
									activekit.Attention(err.Error())
									return nil
								}
								fname := activekit.Promt("Print filename: ")
								if err := ioutil.WriteFile(fname, []byte(data), os.ModePerm); err != nil {
									logrus.WithError(err).Errorf("unable to write service data to file")
									activekit.Attention(err.Error())
									return nil
								}
								fmt.Println("OK")
								return nil
							},
						},
						{
							Label: "Exit",
							Action: func() error {
								if activekit.YesNo("Are you sure? [Y/N]: ") {
									os.Exit(0)
								}
								return nil
							},
						},
					},
				}).Run()
			}
		},
	}
	command.PersistentFlags().
		StringVarP(&createServiceConfig.File, "file", "f", createServiceConfig.File, "file with service data")
	command.PersistentFlags().
		BoolVar(&createServiceConfig.Force, "force", false, "create service without confirmation")
	return command
}
