package clingress

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/ingress/activeingress"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Create(ctx *context.Context) *cobra.Command {
	var flags activeingress.Flags
	command := &cobra.Command{
		Use:     "ingress",
		Aliases: aliases,
		Short:   "create ingress",
		Long:    "Create ingress. Available options: TLS with LetsEncrypt and custom certs.",
		Example: "chkit create ingress [--force] [--filename ingress.json] [-n prettyNamespace]",
		Run: func(cmd *cobra.Command, args []string) {
			flagIngress, err := flags.Ingress()
			if err != nil {
				ferr.Println(err)
				ctx.Exit(1)
			}

			if flags.Force {
				if err := activeingress.ValidateIngress(flagIngress); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				if err := ctx.Client.CreateIngress(ctx.GetNamespace().ID, flagIngress); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Printf("Congratulations! Ingress %s created!\n", flagIngress.Name)
				return
			}

			services, err := ctx.Client.GetServiceList(ctx.GetNamespace().ID)
			services = services.AvailableForIngress()
			if err != nil {
				activekit.Attention(fmt.Sprintf("Unable to get service list!\n%v", err))
				ctx.Exit(1)
			}
			ingr, err := activeingress.Wizard(activeingress.Config{
				Services: services,
				Ingress:  &flagIngress,
			})
			if err != nil {
				activekit.Attention(err.Error())
				ctx.Exit(1)
			}
			fmt.Println(ingr.RenderTable())
			if activekit.YesNo("Are you sure you want create ingress %q?", ingr.Name) {
				if err := ctx.Client.CreateIngress(ctx.GetNamespace().ID, ingr); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Printf("Congratulations! Ingress %s created!\n", ingr.Name)
			}
		},
	}

	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
