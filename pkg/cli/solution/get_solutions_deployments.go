package clisolution

import (
	"fmt"

	"os"

	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var aliases_depl = []string{"sol_depl", "solution_deployments", "sol_deploy", "soldeploy"}

func GetDepl(ctx *context.Context) *cobra.Command {
	var getDeplDataConfig = struct {
		configuration.ExportConfig
	}{
		configuration.ExportConfig{
			Format: configuration.PRETTY,
		},
	}
	command := &cobra.Command{
		Use:     "soldepl",
		Short:   "show deployment data",
		Long:    "Print deployment data.",
		Example: "chkit get solution solution_name [-o yaml/json] [-f output_file]",
		Aliases: aliases_depl,
		Run: func(command *cobra.Command, args []string) {
			deplData, err := func() (model.Renderer, error) {
				if len(args) == 1 {
					logrus.Debugf("getting deployment from %q", ctx.Namespace)
					list, err := ctx.Client.GetSolutionDeployments(ctx.Namespace.ID, args[0])
					if err != nil {
						return nil, err
					}
					return list, nil
				} else {
					command.Help()
					os.Exit(1)
					return nil, nil
				}
			}()
			if err != nil {
				logrus.WithError(err).Errorf("unable to get deployment data")
				fmt.Printf("%v :(\n", err)
				return
			}
			if err := configuration.ExportData(deplData, getDeplDataConfig.ExportConfig); err != nil {
				logrus.WithError(err).Errorf("unable to export data")
				angel.Angel(ctx, err)
			}
		},
	}

	command.PersistentFlags().
		StringVarP((*string)(&getDeplDataConfig.Format), "output", "o", "", "output format (yaml/json)")
	command.PersistentFlags().
		StringVarP(&getDeplDataConfig.Filename, "file", "f", "", "output file")

	return command
}
