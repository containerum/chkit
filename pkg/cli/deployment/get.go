package clideployment

import (
	"fmt"

	"github.com/containerum/chkit/pkg/configuration"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/strset"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	ErrNoNamespaceSpecified chkitErrors.Err = "no namespace specified"
)

var aliases = []string{"depl", "deployments", "deploy"}
var getDeplDataConfig = struct {
	configuration.ExportConfig
}{
	configuration.ExportConfig{
		Format: configuration.PRETTY,
	},
}

var Get = &cobra.Command{
	Use:     "deployment",
	Short:   "shows deployment data",
	Long:    "Shows deployment data",
	Example: "namespace deployment_names... [-n namespace_label]",
	Aliases: aliases,
	Run: func(command *cobra.Command, args []string) {
		deplData, err := func() (model.Renderer, error) {
			switch len(args) {
			case 0:
				logrus.Debugf("getting deployment from %q", context.GlobalContext.Namespace)
				list, err := context.GlobalContext.Client.GetDeploymentList(context.GlobalContext.Namespace)
				if err != nil {
					return nil, err
				}
				return list, nil
			default:
				deplNames := strset.NewSet(args)
				var showList deployment.DeploymentList = make([]deployment.Deployment, 0) // prevents panic
				list, err := context.GlobalContext.Client.GetDeploymentList(context.GlobalContext.Namespace)
				if err != nil {
					return nil, err
				}
				for _, depl := range list {
					if deplNames.Have(depl.Name) {
						showList = append(showList, depl)
					}
				}
				return showList, nil
			}
		}()
		if err != nil {
			logrus.WithError(err).Errorf("unable to get deployment data")
			fmt.Printf("%v :(\n", err)
			return
		}
		if err := configuration.ExportData(deplData, configuration.ExportConfig{}); err != nil {
			logrus.WithError(err).Errorf("unable to export data")
			angel.Angel(err)
		}
	},
}

func init() {
	Get.PersistentFlags().
		StringVarP((*string)(&getDeplDataConfig.Format), "output", "o", "", "output format (yaml/json)")
	Get.PersistentFlags().
		StringVarP(&getDeplDataConfig.Filename, "file", "f", "", "output file")
}
