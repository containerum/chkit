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
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Replace(ctx *context.Context) *cobra.Command {
	var file string
	var force bool
	var flagService service.Service
	var flagPort = service.Port{
		Port: new(int),
	}

	command := &cobra.Command{
		Use:     "service",
		Aliases: aliases,
		Short:   "replace service",
		Long: `Replaces service.
Has an one-line mode, suitable for integration with other tools, and an interactive wizard mode`,
		Run: func(cmd *cobra.Command, args []string) {
			serv := service.Service{}
			if cmd.Flag("file").Changed {
				var err error
				serv, err = servactive.FromFile(file)
				if err != nil {
					logrus.WithError(err).Errorf("unable to load service data from file %s", file)
					fmt.Printf("Unable to load service data from file :(\n%v", err)
					os.Exit(1)
				}
			} else if cmd.Flag("force").Changed {
				serv = flagService
			}
			if cmd.Flag("force").Changed {
				if len(args) != 1 {
					cmd.Help()
					return
				}
				serv.Name = args[0]
				serv.Ports = []service.Port{flagPort}

				oldServ, err := ctx.Client.GetService(ctx.Namespace, args[0])
				if err != nil {
					activekit.Attention(err.Error())
					os.Exit(1)
				}
				if !cmd.Flag("port").Changed {
					flagPort.Port = nil
				}
				if !cmd.Flag("port-name").Changed {
					serv.Ports = oldServ.Ports
				}
				if !cmd.Flag("deployment").Changed {
					serv.Deploy = oldServ.Deploy
				}
				if !cmd.Flag("domain").Changed {
					serv.Domain = oldServ.Domain
				}
				if err := servactive.ValidateService(serv); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(serv.RenderTable())
				if err := ctx.Client.ReplaceService(ctx.Namespace, serv); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println("OK")
				return
			} else {
				if len(args) == 0 {
					list, err := ctx.Client.GetServiceList(ctx.Namespace)
					if err != nil {
						activekit.Attention(err.Error())
						os.Exit(1)
					}
					var menu []*activekit.MenuItem
					for _, s := range list {
						menu = append(menu, &activekit.MenuItem{
							Label: s.Name,
							Action: func(d service.Service) func() error {
								return func() error {
									serv = d
									return nil
								}
							}(s),
						})
					}
					(&activekit.Menu{
						Title: "Choose service to replace",
						Items: menu,
					}).Run()
				} else {
					var err error
					serv, err = ctx.Client.GetService(ctx.Namespace, args[0])
					if err != nil {
						activekit.Attention(err.Error())
						os.Exit(1)
					}
				}
			}
			serv, err := servactive.ReplaceWizard(servactive.ConstructorConfig{
				Service: &serv,
			})
			if err != nil {
				logrus.WithError(err).Errorf("unable to replace service")
				fmt.Println(err)
				os.Exit(1)
			}
			for {
				_, err := (&activekit.Menu{
					Items: []*activekit.MenuItem{
						{
							Label: fmt.Sprintf("Update service %q on server", serv.Name),
							Action: func() error {
								fmt.Println(serv.RenderTable())
								if activekit.YesNo(fmt.Sprintf("Are you sure you want to update service %q on server?", serv.Name)) {
									err := ctx.Client.ReplaceService(ctx.Namespace, serv)
									if err != nil {
										logrus.WithError(err).Errorf("unable to replace service %q", serv.Name)
										fmt.Println(err)
										return nil
									}
									fmt.Printf("Congratulations! Service %q updated!\n", serv.Name)
								}
								return nil
							},
						},
						{
							Label: "Edit service",
							Action: func() error {
								var err error
								serv, err = servactive.ReplaceWizard(servactive.ConstructorConfig{
									Service: &serv,
								})
								if err != nil {
									logrus.WithError(err).Errorf("unable to update service")
									fmt.Println(err)
									os.Exit(1)
								}
								return nil
							},
						},
						{
							Label: "Print to terminal",
							Action: activekit.ActionWithErr(func() error {
								if data, err := serv.RenderYAML(); err != nil {
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
								data, err := serv.RenderJSON()
								if err != nil {
									return err
								}
								if err := ioutil.WriteFile(filename, []byte(data), os.ModePerm); err != nil {
									logrus.WithError(err).Errorf("unable to save service %q to file", serv.Name)
									fmt.Printf("Unable to save service to file :(\n%v", err)
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
		StringVar(&file, "file", "", "create service from file")
	command.PersistentFlags().
		BoolVarP(&force, "force", "f", false, "suppress confirmation")

	command.PersistentFlags().
		StringVar(&flagService.Deploy, "deployment", "", "deployment name, optional")
	command.PersistentFlags().
		StringVar(&flagService.Domain, "domain", "", "service domain, optional")
	command.PersistentFlags().
		IntVar(flagPort.Port, "port", 0, "service external port, optional")
	command.PersistentFlags().
		IntVar(&flagPort.TargetPort, "target-port", 80, "service target port, optional")
	command.PersistentFlags().
		StringVar(&flagPort.Name, "port-name", "", "service port name")
	command.PersistentFlags().
		StringVar(&flagPort.Protocol, "protocol", "TCP", "service port protocol, optional")
	return command
}
