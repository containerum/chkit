package cliserv

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/export"
	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/model/service/servactive"
	"github.com/containerum/chkit/pkg/porta"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Replace(ctx *context.Context) *cobra.Command {
	var flags struct {
		servactive.Flags
		porta.Importer
		porta.Exporter
	}
	exportConfig := export.ExportConfig{}
	command := &cobra.Command{
		Use:     "service",
		Aliases: aliases,
		Short:   "Replace service.",
		Long: `Replace service.\n` +
			`Runs in one-line mode, suitable for integration with other tools, and in interactive wizard mode.`,
		Run: func(cmd *cobra.Command, args []string) {
			var external bool
			var logger = coblog.Logger(cmd)
			logger.Struct(flags)
			logger.Debugf("running replace service command")
			var flagSvc service.Service
			if flags.ImportActivated() {
				if err := flags.Import(&flagSvc); err != nil {
					ferr.Printf("unable to import service:\n%v\n", err)
					ctx.Exit(1)
				}
			} else {
				var err error
				flagSvc, err = flags.Service()
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
			}
			var svc service.Service
			if flags.Force {
				if len(args) != 1 {
					cmd.Help()
					return
				}
				oldServ, err := ctx.Client.GetService(ctx.GetNamespace().ID, args[0])
				if err != nil {
					activekit.Attention(err.Error())
					ctx.Exit(1)
				}
				if oldServ.Domain != "" {
					external = true
				}
				if len(flagSvc.Ports) != 0 {
					oldServ.Ports = append(oldServ.Ports, flagSvc.Ports[0])
				}
				if flagSvc.Deploy != "" {
					oldServ.Deploy = flagSvc.Deploy
				}
				if err := servactive.ValidateService(oldServ); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				if flags.ExporterActivated() {
					if err := flags.Export(oldServ); err != nil {
						ferr.Printf("unable to export service:\n%v\n", err)
						ctx.Exit(1)
					}
					return
				}
				if err := ctx.Client.ReplaceService(ctx.GetNamespace().ID, oldServ); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Printf("Congratulations! Service %s updated!\n", oldServ.Name)
				return
			} else {
				if len(args) == 0 {
					list, err := ctx.Client.GetServiceList(ctx.GetNamespace().ID)
					if err != nil {
						activekit.Attention(err.Error())
						ctx.Exit(1)
					}
					var menu []*activekit.MenuItem
					for _, s := range list {
						menu = append(menu, &activekit.MenuItem{
							Label: s.Name,
							Action: func(d service.Service) func() error {
								return func() error {
									svc = d
									if svc.Domain != "" {
										external = true
									}
									return nil
								}
							}(s),
						})
					}
					if err := export.ExportData(list, exportConfig); err != nil {
						logrus.WithError(err).Errorf("unable to export data")
						angel.Angel(ctx, err)
					}
					(&activekit.Menu{
						Title: "Choose service to replace",
						Items: menu,
					}).Run()
				} else {
					var err error
					svc, err = ctx.Client.GetService(ctx.GetNamespace().ID, args[0])
					if svc.Domain != "" {
						external = true
					}
					if err != nil {
						activekit.Attention(err.Error())
						ctx.Exit(1)
					}
				}
			}
			depList, err := ctx.Client.GetDeploymentList(ctx.GetNamespace().ID)
			if err != nil {
				logrus.WithError(err).Errorf("unable to get deployment list")
				fmt.Println("Unable to get deployment list :(")
			}
			if len(flagSvc.Ports) != 0 {
				svc.Ports = append(svc.Ports, flagSvc.Ports[0])
			}
			if flagSvc.Deploy != "" {
				svc.Deploy = flagSvc.Deploy
			}
			svc, err = servactive.ReplaceWizard(servactive.ConstructorConfig{
				External:    external,
				Deployments: depList.Names(),
				Service:     &svc,
			})
			if err != nil {
				activekit.Attention(err.Error())
				ctx.Exit(1)
			}
			if activekit.YesNo("Are you sure you want update service %q?", svc.Name) {
				if err := ctx.Client.ReplaceService(ctx.GetNamespace().ID, svc); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Printf("Congratulations! Service %s updated!\n", svc.Name)
			}
			svc, err = ctx.Client.GetService(ctx.GetNamespace().ID, svc.Name)
			if err != nil {
				logrus.WithError(err).Errorf("unable to get service")
				fmt.Println("Unable to get service :(")
			}
			fmt.Println(svc.RenderTable())
			(&activekit.Menu{
				Items: activekit.MenuItems{
					{
						Label: "Edit service " + svc.Name,
						Action: func() error {
							changedSvc, err := servactive.ReplaceWizard(servactive.ConstructorConfig{
								External:    external,
								Deployments: depList.Names(),
								Service:     &svc,
							})
							if err != nil {
								ferr.Println(err)
								return nil
							}
							svc = changedSvc
							if activekit.YesNo("Push changes to server?") {
								if err := ctx.Client.ReplaceService(ctx.GetNamespace().ID, svc); err != nil {
									ferr.Printf("unable to update service on server:\n%v\n", err)
								}
							}
							return nil
						},
					},
					{
						Label: "Export service to file",
						Action: func() error {
							var fname = activekit.Promt("Type filename: ")
							fname = strings.TrimSpace(fname)
							if fname != "" {
								if err := (porta.Exporter{OutFile: fname}.Export(svc)); err != nil {
									ferr.Printf("unable to export service:\n%v\n", err)
								}
							}
							return nil
						},
					},
				},
			}).Run()
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
