package cliserv

import (
	"fmt"

	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/strset"
	"github.com/spf13/cobra"
)

var aliases = []string{"srv", "services", "svc", "serv"}

var getServiceConfig = struct {
	configuration.ExportConfig
}{}

var Get = &cobra.Command{
	Use:     "service",
	Aliases: aliases,
	Short:   "shows service info",
	Long:    "chkit get service service_label [-o yaml/json] [-f output_file]",
	Example: "Shows service info",
	Run: func(cmd *cobra.Command, args []string) {
		serviceData, err := func() (model.Renderer, error) {
			switch len(args) {
			case 0:
				list, err := context.GlobalContext.Client.GetServiceList(context.GlobalContext.Namespace)
				return list, err
			case 1:
				svc, err := context.GlobalContext.Client.GetDeployment(context.GlobalContext.Namespace, args[0])
				return svc, err
			default:
				list, err := context.GlobalContext.Client.GetServiceList(context.GlobalContext.Namespace)
				var filteredList service.ServiceList
				names := strset.NewSet(args)
				for _, svc := range list {
					if names.Have(svc.Name) {
						filteredList = append(filteredList, svc)
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

func init() {
	Get.PersistentFlags().
		StringVarP((*string)(&getServiceConfig.Format), "output", "o", "", "output format [yaml/json]")
	Get.PersistentFlags().
		StringVarP(&getServiceConfig.Filename, "file", "f", "-", "output file")
}
