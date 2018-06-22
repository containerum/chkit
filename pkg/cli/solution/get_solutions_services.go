package clisolution

import (
	"fmt"

	"os"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/export"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var aliases_svc = []string{"sol_svc", "solution_svc", "solution_services"}

func GetSvc(ctx *context.Context) *cobra.Command {
	var getDeplDataConfig = struct {
		export.ExportConfig
	}{
		export.ExportConfig{
			Format: export.PRETTY,
		},
	}
	command := &cobra.Command{
		Use:     "solsvc",
		Short:   "Show solution services data",
		Example: "chkit get solution_svc solution_name [-o yaml/json] [-f output_file]",
		Aliases: aliases_svc,
		Run: func(command *cobra.Command, args []string) {
			deplData, err := func() (model.Renderer, error) {
				if len(args) == 1 {
					logrus.Debugf("getting services from %q", ctx.Namespace)
					list, err := ctx.Client.GetSolutionServices(ctx.Namespace.ID, args[0])
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
			if err := export.ExportData(deplData, getDeplDataConfig.ExportConfig); err != nil {
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
