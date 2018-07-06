package clingress

import (
	"fmt"

	"strings"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/export"
	"github.com/containerum/chkit/pkg/model/ingress"
	"github.com/containerum/chkit/pkg/model/ingress/activeingress"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Replace(ctx *context.Context) *cobra.Command {
	var flags activeingress.Flags
	exportConfig := export.ExportConfig{}
	command := &cobra.Command{
		Use:     "ingress",
		Short:   "Replace ingress.",
		Aliases: aliases,
		Long: "Replace ingress with a new one, use --force flag to write one-liner command, " +
			"omitted attributes are inherited from the previous ingress.",
		Example: "chkit replace ingress $INGRESS [--force] [--service $SERVICE] [--port 80] [--tls-secret letsencrypt]",
		Run: func(cmd *cobra.Command, args []string) {
			logger := coblog.Logger(cmd)
			logger.Struct(flags)
			logger.Debugf("running replace ingress command")
			var ingr ingress.Ingress
			if len(args) == 1 {
				ingrName := args[0]
				var err error
				ingr, err = ctx.Client.GetIngress(ctx.GetNamespace().ID, ingrName)
				if err != nil {
					logger.WithError(err).Errorf("unable to get previous ingress")
					activekit.Attention("unable to get previous ingress:\n%v", err)
					ctx.Exit(1)
				}
			} else if !flags.Force {
				ingrList, err := ctx.Client.GetIngressList(ctx.GetNamespace().ID)
				if err != nil {
					logger.WithError(err).Errorf("unable to get ingress list")
					activekit.Attention("Unable to get ingress list:\n%v", err)
					ctx.Exit(1)
				}
				if ingrList.Len() == 0 {
					logger.Errorf("no ingresses exists")
					fmt.Println("no ingresses exists\n")
					ctx.Exit(1)
				} else if ingrList.Len() == 1 {
					ingr = ingrList.Head()
				} else if ingrList.Len() > 1 {
					if err := export.ExportData(ingrList, exportConfig); err != nil {
						logrus.WithError(err).Errorf("unable to export data")
						angel.Angel(ctx, err)
					}
					var menu activekit.MenuItems
					for _, i := range ingrList {
						menu = menu.Append(&activekit.MenuItem{
							Label: i.Name,
							Action: func(i ingress.Ingress) func() error {
								return func() error {
									ingr = i.Copy()
									return nil
								}
							}(i),
						})
					}
					(&activekit.Menu{
						Title: "Select ingress to replace",
						Items: menu,
					}).Run()
				}
			} else {
				cmd.Help()
				ctx.Exit(1)
			}
			ingrChanged, err := flags.Ingress()
			if err != nil {
				logger.WithError(err).Errorf("unable to load ingress from flags")
				activekit.Attention("Unable to load ingress from flags:\n%v", err)
				ctx.Exit(1)
			}

			if ingrChanged.Rules != nil {
				if ingrChanged.Rules[0].TLSSecret != "" {
					ingr.Rules[0].TLSSecret = ingrChanged.Rules[0].TLSSecret
				}
				if ingrChanged.Rules[0].Paths != nil {
					ingr.Rules[0].Paths = ingrChanged.Rules[0].Paths
				}
				ingr.Rules[0].Host = strings.TrimRight(ingr.Rules[0].Host, ".hub.containerum.io")
			}
			if flags.Force {
				if err := activeingress.ValidateIngress(ingr); err != nil {
					logger.WithError(err).Errorf("invalid flag-defined ingress")
					activekit.Attention("Invalid ingress:\n%v", err)
					ctx.Exit(1)
				}
				if err := ctx.Client.ReplaceIngress(ctx.GetNamespace().ID, ingr); err != nil {
					logger.WithError(err).Errorf("unable to replace ingress")
					activekit.Attention("Unable to replace ingress %q:\n%v", ingr.Name, err)
					ctx.Exit(1)
				}
				fmt.Println("OK")
				return
			}
			services, err := ctx.Client.GetServiceList(ctx.GetNamespace().ID)
			if err != nil {
				logger.WithError(err).Errorf("unable to get service list")
				activekit.Attention("Unable to get service list:\n%v", err)
				ctx.Exit(1)
			}
			services = services.AvailableForIngress()
			ingr.Rules[0].Host = strings.TrimRight(ingr.Rules[0].Host, ".hub.containerum.io")
			ingr, err = activeingress.EditWizard(activeingress.Config{
				Services: services,
				Ingress:  &ingr,
			})
			if err != nil {
				logger.WithError(err).Errorf("unable to build ingress")
				activekit.Attention("Unable to build ingress:\n%v", err)
			}
			if activekit.YesNo("Do you really want to replace ingress %q?", ingr.Name) {
				if err := activeingress.ValidateIngress(ingr); err != nil {
					logger.WithError(err).Errorf("invalid flag-defined ingress")
					activekit.Attention("Invalid ingress:\n%v", err)
					ctx.Exit(1)
				}
				if err := ctx.Client.ReplaceIngress(ctx.GetNamespace().ID, ingr); err != nil {
					logger.WithError(err).Errorf("unable to replace ingress")
					activekit.Attention("Unable to replace ingress %q:\n%v", ingr.Name, err)
					ctx.Exit(1)
				}
				fmt.Println("OK")
			}
		},
	}

	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
