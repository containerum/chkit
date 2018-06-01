package clingress

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/ingress"
	"github.com/containerum/chkit/pkg/model/ingress/activeingress"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/spf13/cobra"
)

func Replace(ctx *context.Context) *cobra.Command {
	var flags = struct {
		Force bool
	}{}
	command := &cobra.Command{
		Use:     "ingress",
		Short:   "patch ingress with new attributes",
		Aliases: aliases,
		Long:    "Replaces ingress with new, use --force flag to write one-liner command, omitted attributes are inherited from previous ingress.",
		Example: "chkit replace ingress $INGRESS [--force] [--service $SERVICE] [--port 80] [--tls-secret letsencrypt]",
		Run: func(cmd *cobra.Command, args []string) {
			logger := coblog.Logger(cmd)
			logger.Debugf("running replace ingress command")
			var ingr ingress.Ingress
			if len(args) == 1 {
				ingrName := args[0]
				var err error
				ingr, err = ctx.Client.GetIngress(ctx.Namespace.ID, ingrName)
				if err != nil {
					logger.WithError(err).Errorf("unable to get previous ingress")
					activekit.Attention("unable to get previous ingress:\n%v", err)
					os.Exit(1)
				}
			} else if !flags.Force {
				ingrList, err := ctx.Client.GetIngressList(ctx.Namespace.ID)
				if err != nil {
					logger.WithError(err).Errorf("unable to get ingress list")
					activekit.Attention("Unable to get ingress list:\n%v", err)
					os.Exit(1)
				}
				fmt.Println(ingrList)
				if ingrList.Len() == 1 {
					ingr = ingrList.Head()
				} else if ingrList.Len() == 0 {
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
				os.Exit(1)
			}
			ingr, ingrChanged := buildIngress(cmd, ingr)
			if cmd.Flag("force").Changed && ingrChanged {
				if err := activeingress.ValidateIngress(ingr); err != nil {
					logger.WithError(err).Errorf("invalid flag-defined ingress")
					activekit.Attention("Invalid ingress:\n%v", err)
					os.Exit(1)
				}
				if err := ctx.Client.ReplaceIngress(ctx.Namespace.ID, ingr); err != nil {
					logger.WithError(err).Errorf("unable to replace ingress")
					activekit.Attention("Unable to replace ingress %q:\n%v", ingr.Name, err)
					os.Exit(1)
				}
				fmt.Println("OK")
				return
			} else if !ingrChanged {
				fmt.Println("Nothing to do")
			}
			services, err := ctx.Client.GetServiceList(ctx.Namespace.ID)
			if err != nil {
				logger.WithError(err).Errorf("unable to get service list")
				activekit.Attention("Unable to get service list:\n%v", err)
				os.Exit(1)
			}
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
					os.Exit(1)
				}
				if err := ctx.Client.ReplaceIngress(ctx.Namespace.ID, ingr); err != nil {
					logger.WithError(err).Errorf("unable to replace ingress")
					activekit.Attention("Unable to replace ingress %q:\n%v", ingr.Name, err)
					os.Exit(1)
				}
				fmt.Println("OK")
			}
		},
	}

	command.PersistentFlags().
		BoolVarP(&flags.Force, "force", "f", false, "replace ingress without confirmation")
	command.PersistentFlags().
		String("host", "", "ingress host, optional")
	command.PersistentFlags().
		String("tls-secret", "", "ingress tls-secret, use 'letsencrypt' for automatic HTTPS, '-' to use HTTP, optional")
	command.PersistentFlags().
		Int("port", 8080, "ingress endpoint port, optional")
	command.PersistentFlags().
		String("service", "", "ingress endpoin service, optional")
	return command
}
