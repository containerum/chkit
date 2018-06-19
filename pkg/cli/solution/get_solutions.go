package clisolution

import (
	"fmt"

	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/solution"
	"github.com/containerum/chkit/pkg/util/strset"
	"github.com/spf13/cobra"
)

func Get(ctx *context.Context) *cobra.Command {
	var getServiceConfig = struct {
		configuration.ExportConfig
	}{}
	command := &cobra.Command{
		Use:     "solution",
		Aliases: aliases,
		Short:   "show solution info",
		Long:    "Show solution info.",
		Example: "chkit get solution solution_name [-o yaml/json] [-f output_file]",
		Run: func(cmd *cobra.Command, args []string) {
			serviceData, err := func() (model.Renderer, error) {
				switch len(args) {
				case 0:
					list, err := ctx.Client.GetRunningSolutionsList(ctx.Namespace.ID)
					return list, err
				case 1:
					sol, err := ctx.Client.GetRunningSolution(ctx.Namespace.ID, args[0])
					return sol, err
				default:
					list, err := ctx.Client.GetRunningSolutionsList(ctx.Namespace.ID)
					var filteredList solution.SolutionsList
					names := strset.NewSet(args)
					for _, sol := range list.Solutions {
						if names.Have(sol.Name) {
							filteredList.Solutions = append(filteredList.Solutions, sol)
						}
					}
					return filteredList, err
				}
			}()
			if err != nil {
				fmt.Println(err)
				return
			}
			if err := configuration.ExportData(serviceData, getServiceConfig.ExportConfig); err != nil {
				fmt.Println(err)
				return
			}

		},
	}
	command.PersistentFlags().
		StringVarP((*string)(&getServiceConfig.Format), "output", "o", "", "output format [yaml/json]")
	command.PersistentFlags().
		StringVarP(&getServiceConfig.Filename, "file", "f", "-", "output file")
	return command
}
