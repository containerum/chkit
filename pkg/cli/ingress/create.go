package clingress

import (
	"fmt"
	"io/ioutil"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/ingress"
	"github.com/containerum/chkit/pkg/model/ingress/activeingress"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/containerum/chkit/pkg/util/host2dnslabel"
	"github.com/spf13/cobra"
)

func Create(ctx *context.Context) *cobra.Command {
	var force bool
	var flagIngress ingress.Ingress
	var flagRule = ingress.Rule{
		TLSSecret: new(string),
	}
	var flagPath ingress.Path
	var tlsSecretFile string

	command := &cobra.Command{
		Use:     "ingress",
		Aliases: aliases,
		Short:   "create ingress",
		Long:    "Create ingress. Available options: TLS with LetsEncrypt and custom certs.",
		Example: "chkit create ingress [--force] [--filename ingress.json] [-n prettyNamespace]",
		Run: func(cmd *cobra.Command, args []string) {
			if !cmd.Flag("tls-secret").Changed {
				flagRule.TLSSecret = nil
			}
			if cmd.Flag("tls-cert").Changed {
				cert, err := ioutil.ReadFile(tlsSecretFile)
				if err != nil {
					fmt.Printf("unable to read cert file: %v\n", err.Error())
					ctx.Exit(1)
				}
				c := string(cert)
				flagRule.TLSSecret = &c
			}
			if cmd.Flag("path").Changed ||
				cmd.Flag("service").Changed ||
				cmd.Flag("port").Changed {
				flagRule.Paths = ingress.PathList{flagPath}
			}
			if cmd.Flag("host").Changed ||
				cmd.Flag("tls-secret").Changed ||
				cmd.Flag("tls-cert").Changed ||
				cmd.Flag("path").Changed ||
				cmd.Flag("service").Changed ||
				cmd.Flag("port").Changed {
				flagIngress.Rules = ingress.RuleList{flagRule}
				flagIngress.Name = host2dnslabel.Host2DNSLabel(flagRule.Host)
			}

			if cmd.Flag("force").Changed {
				if err := activeingress.ValidateIngress(flagIngress); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				if err := ctx.Client.CreateIngress(ctx.GetNamespace().ID, flagIngress); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
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

	command.PersistentFlags().
		BoolVarP(&force, "force", "f", false, "create ingress without confirmation")
	command.PersistentFlags().
		StringVar(&flagRule.Host, "host", "", "ingress host (example: prettyblog.io), required")
	command.PersistentFlags().
		StringVar(flagRule.TLSSecret, "tls-secret", "", "TLS secret string, optional")
	command.PersistentFlags().
		StringVar(&tlsSecretFile, "tls-cert", "", "TLS cert file, optional")
	command.PersistentFlags().
		StringVar(&flagPath.Path, "path", "", "path to endpoint (example: /content/pages), optional")
	command.PersistentFlags().
		StringVar(&flagPath.ServiceName, "service", "", "ingress endpoint service, required")
	command.PersistentFlags().
		IntVar(&flagPath.ServicePort, "port", 8080, "ingress endpoint port (example: 80, 443), optional")
	return command
}
