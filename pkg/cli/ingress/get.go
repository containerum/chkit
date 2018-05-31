package clingress

import (
	"fmt"

	"os"

	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/ingress"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/strset"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var aliases = []string{"ingr", "ingresses", "ing"}

func Get(ctx *context.Context) *cobra.Command {
	exportConfig := configuration.ExportConfig{}
	command := &cobra.Command{
		Use:     "ingress",
		Short:   "show ingress data",
		Long:    "Shows ingress data",
		Example: "chkit get ingress ingress_names... [-n namespace_label] [-o yaml/json]",
		Aliases: aliases,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := prerun.PreRun(ctx); err != nil {
				angel.Angel(ctx, err)
				os.Exit(1)
			}
			if cmd.Flags().Changed("namespace") {
				ctx.Namespace.ID, _ = cmd.Flags().GetString("namespace")
			}
		},
		Run: func(command *cobra.Command, args []string) {
			ingrData, err := func() (model.Renderer, error) {
				switch len(args) {
				case 0:
					logrus.Debugf("getting ingress from %q", ctx.Namespace)
					list, err := ctx.Client.GetIngressList(ctx.Namespace.ID)
					if err != nil {
						return nil, err
					}
					return list, nil
				case 1:
					logrus.Debugf("getting ingress from %q", ctx.Namespace)
					ingr, err := ctx.Client.GetIngress(ctx.Namespace.ID, args[0])
					if err != nil {
						return nil, err
					}
					return ingr, nil
				default:
					deplNames := strset.NewSet(args)
					var showList = make(ingress.IngressList, 0) // prevents panic
					list, err := ctx.Client.GetIngressList(ctx.Namespace.ID)
					if err != nil {
						return nil, err
					}
					for _, ingr := range list {
						if deplNames.Have(ingr.Host()) {
							showList = append(showList, ingr)
						}
					}
					return showList, nil
				}
			}()
			if err != nil {
				logrus.WithError(err).Errorf("unable to get ingress data")
				fmt.Printf("%v :(\n", err)
				return
			}
			if err := configuration.ExportData(ingrData, exportConfig); err != nil {
				logrus.WithError(err).Errorf("unable to export data")
				angel.Angel(ctx, err)
			}
		},
	}

	command.PersistentFlags().
		StringVarP((*string)(&exportConfig.Format), "output", "o", "", "output format (yaml/json)")
	command.PersistentFlags().
		StringVarP(&exportConfig.Filename, "file", "f", "", "output file")

	return command
}
