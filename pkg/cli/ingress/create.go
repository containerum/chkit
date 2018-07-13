package clingress

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/ingress"
	"github.com/containerum/chkit/pkg/model/ingress/activeingress"
	"github.com/containerum/chkit/pkg/porta"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Create(ctx *context.Context) *cobra.Command {
	var flags struct {
		activeingress.Flags
		porta.Importer
		porta.Exporter
	}
	command := &cobra.Command{
		Use:     "ingress",
		Aliases: aliases,
		Short:   "create ingress",
		Long:    "Create ingress. Available options: TLS with LetsEncrypt and custom certs.",
		Example: "chkit create ingress [--force] [--filename ingress.json] [-p project]",
		Run: func(cmd *cobra.Command, args []string) {
			var logger = coblog.Logger(cmd)
			logger.Struct(flags)
			logger.Debugf("running create ingress command")
			var ingr ingress.Ingress
			if flags.ImportActivated() {
				if err := flags.Import(&ingr); err != nil {
					ferr.Printf("unable to import ingress:\n%v\n", err)
					ctx.Exit(1)
				}
			} else {
				var err error
				ingr, err = flags.Ingress()
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
			}
			if flags.Force {
				if err := activeingress.ValidateIngress(ingr); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				if flags.ExporterActivated() {
					if err := flags.Export(ingr); err != nil {
						ferr.Printf("unable to export ingress:\n%v\n", err)
						ctx.Exit(1)
					}
					return
				}
				if err := ctx.Client.CreateIngress(ctx.GetNamespace().ID, ingr); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Printf("Congratulations! Ingress %s created!\n", ingr.Name)
				return
			}
			services, err := ctx.Client.GetServiceList(ctx.GetNamespace().ID)
			services = services.AvailableForIngress()
			if err != nil {
				activekit.Attention(fmt.Sprintf("Unable to get service list!\n%v", err))
				ctx.Exit(1)
			}
			ingr, err = activeingress.Wizard(activeingress.Config{
				Services: services,
				Ingress:  &ingr,
			})
			if err != nil {
				activekit.Attention(err.Error())
				ctx.Exit(1)
			}
			if activekit.YesNo("Are you sure you want to create ingress %q?", ingr.Name) {
				if err := ctx.Client.CreateIngress(ctx.GetNamespace().ID, ingr); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Printf("Congratulations! Ingress %s created!\n", ingr.Name)
			} else {
				ctx.Exit(0)
			}
			fmt.Println(ingr.RenderTable())
			(&activekit.Menu{
				Items: activekit.MenuItems{
					{
						Label: "Edit ingress " + ingr.Name,
						Action: func() error {
							changedIng, err := activeingress.Wizard(activeingress.Config{
								Services: services,
								Ingress:  &ingr,
							})
							if err != nil {
								ferr.Println(err)
								return nil
							}
							ingr = changedIng
							if activekit.YesNo("Push changes to server?") {
								if err := ctx.Client.ReplaceIngress(ctx.GetNamespace().ID, ingr); err != nil {
									ferr.Printf("unable to update ingress on server:\n%v\n", err)
								}
							}
							return nil
						},
					},
					{
						Label: "Exit",
						Action: func() error {
							ctx.Exit(0)
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
