package cliserv

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/model/service/servactive"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Create(ctx *context.Context) *cobra.Command {
	var createServiceConfig = struct {
		File          string
		Force         bool
		FlagService   service.Service
		FlagPort      service.Port
		FlagServiceIP string
	}{
		File: "-",
		FlagPort: service.Port{
			Port: new(int),
		},
	}
	command := &cobra.Command{
		Use:     "service",
		Aliases: aliases,
		Short:   "create service",
		Long:    "create service for provided pod in provided namespace",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.WithField("command", "create serv").Debugf("start serv creation")
			depList, err := ctx.Client.GetDeploymentList(ctx.Namespace)
			if err != nil {
				logrus.WithError(err).Errorf("unable to get deployment list")
				fmt.Println("Unable to get deployment list :(")
			}
			wizardConfig := servactive.ConstructorConfig{
				Deployments: depList.Names(),
			}
			if cmd.Flag("file").Changed && createServiceConfig.File != "" && createServiceConfig.File != "-" {
				logrus.Debugf("loading serv from file")
				serv, err := servactive.FromFile(createServiceConfig.File)
				if err != nil {
					logrus.WithError(err).Errorf("unable to load serv from file")
					activekit.Attention(err.Error())
					os.Exit(1)
				}
				if createServiceConfig.Force {
					if err := ctx.Client.CreateService(ctx.Namespace, serv); err != nil {
						logrus.WithError(err).Errorf("unable to create serv %q in namespace %q", serv.Name, ctx.Namespace)
						activekit.Attention(err.Error())
						os.Exit(1)
					}
					fmt.Printf("Service %q created\n", serv.Name)
					return
				}
				wizardConfig.Service = &serv
			} else if createServiceConfig.Force {
				if !cmd.Flag("port").Changed {
					createServiceConfig.FlagPort.Port = nil
				}
				createServiceConfig.FlagService.Ports = []service.Port{createServiceConfig.FlagPort}
				if err := servactive.ValidateService(createServiceConfig.FlagService); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if err := ctx.Client.CreateService(ctx.Namespace, createServiceConfig.FlagService); err != nil {
					logrus.WithError(err).Errorf("unable to create serv %q in namespace %q", createServiceConfig.FlagService.Name, ctx.Namespace)
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println("OK")
				return
			}
			serv, err := servactive.Wizard(wizardConfig)
			if err != nil {
				logrus.WithError(err).Errorf("unable to create serv")
				fmt.Println("Unable to create serv :(")
				os.Exit(1)
			}
			fmt.Println(serv.RenderTable())
			for {
				(&activekit.Menu{
					Items: []*activekit.MenuItem{
						{
							Label: "Push serv to server",
							Action: func() error {
								if activekit.YesNo("Are you sure?") {
									if err := ctx.Client.CreateService(ctx.Namespace, serv); err != nil {
										logrus.WithError(err).Errorf("unable to create serv %q in namespace %q", serv.Name, ctx.Namespace)
										activekit.Attention(err.Error())
										return nil
									}
								}
								logrus.WithError(err).Errorf("serv %q in namespace %q created", serv.Name, ctx.Namespace)
								fmt.Printf("Service %q created\n", serv.Name)
								return nil
							},
						},
						{
							Label: "Edit serv",
							Action: func() error {
								s, err := servactive.Wizard(servactive.ConstructorConfig{
									Service:     &serv,
									Deployments: depList.Names(),
								})
								if err != nil {
									logrus.WithError(err).Errorf("error while interactive serv creation")
									activekit.Attention(err.Error())
									os.Exit(1)
								}
								serv = s
								return nil
							},
						},
						{
							Label: "Print to terminal",
							Action: func() error {
								data, err := serv.RenderYAML()
								if err != nil {
									logrus.WithError(err).Errorf("unable to render serv to yaml")
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
								logrus.Debugf("saving serv to file")
								data, err := serv.RenderJSON()
								if err != nil {
									logrus.WithError(err).Errorf("unable to render serv to json")
									activekit.Attention(err.Error())
									return nil
								}
								fname := activekit.Promt("Print filename: ")
								if err := ioutil.WriteFile(fname, []byte(data), os.ModePerm); err != nil {
									logrus.WithError(err).Errorf("unable to write serv data to file")
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
								if activekit.YesNo("Are you sure?") {
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

	command.PersistentFlags().
		StringVar(&createServiceConfig.FlagService.Name, "name", namegen.ColoredPhysics(), "service name, optional")
	command.PersistentFlags().
		StringVar(&createServiceConfig.FlagService.Deploy, "deploy", "", "service deployment, required")
	command.PersistentFlags().
		IntVar(&createServiceConfig.FlagPort.TargetPort, "target-port", 80, "service target port, optional")
	command.PersistentFlags().
		IntVar(createServiceConfig.FlagPort.Port, "port", 0, "service port, optional")
	command.PersistentFlags().
		StringVar(&createServiceConfig.FlagPort.Protocol, "proto", "TCP", "service protocol, optional")
	command.PersistentFlags().
		StringVar(&createServiceConfig.FlagPort.Name, "port-name", namegen.Aster()+"-"+namegen.Color(), "service port name, optional")
	return command
}
