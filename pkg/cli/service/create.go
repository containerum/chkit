package cliserv

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/model/service/servactive"
	"github.com/containerum/chkit/pkg/porta"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Create(ctx *context.Context) *cobra.Command {
	var flags struct {
		servactive.Flags
		porta.Importer
		porta.Exporter
	}
	command := &cobra.Command{
		Use:     "service",
		Aliases: aliases,
		Short:   "create service",
		Long:    "Create service for the specified pod in the specified namespace.",
		Run: func(cmd *cobra.Command, args []string) {
			var logger = coblog.Logger(cmd)
			logger.Struct(flags)
			logger.Debugf("running create service command")
			var svc service.Service
			if flags.ImportActivated() {
				if err := flags.Import(&svc); err != nil {
					ferr.Printf("unable to import service:\n%v\n", err)
					ctx.Exit(1)
				}
			} else {
				var err error
				svc, err = flags.Service()
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
			}
			if flags.Force {
				if err := servactive.ValidateService(svc); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				if flags.ExporterActivated() {
					if err := flags.Export(svc); err != nil {
						ferr.Printf("unable to export service:\n%v\n", err)
						ctx.Exit(1)
					}
					return
				}
				if err := ctx.Client.CreateService(ctx.GetNamespace().ID, svc); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Printf("Congratulations! Service %s created!\n", svc.Name)
				return
			}
			depList, err := ctx.Client.GetDeploymentList(ctx.GetNamespace().ID)
			if err != nil {
				logrus.WithError(err).Errorf("unable to get deployment list")
				fmt.Println("Unable to get deployment list :(")
			}
			svc, err = servactive.Wizard(servactive.ConstructorConfig{
				Deployments: depList.Names(),
				Service:     &svc,
			})
			if err != nil {
				activekit.Attention(err.Error())
				ctx.Exit(1)
			}
			if activekit.YesNo("Are you sure you want create service %q?", svc.Name) {
				if err := ctx.Client.CreateService(ctx.GetNamespace().ID, svc); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Printf("Congratulations! Service %s created!\n", svc.Name)
			}
			fmt.Println(svc.RenderTable())
			(&activekit.Menu{
				Items: activekit.MenuItems{
					{
						Label: "Edit service " + svc.Name,
						Action: func() error {
							changedSvc, err := servactive.ReplaceWizard(servactive.ConstructorConfig{
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
