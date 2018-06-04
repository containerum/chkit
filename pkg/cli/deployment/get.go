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

func Get(ctx *context.Context) *cobra.Command {
	var getDeplDataConfig = struct {
		configuration.ExportConfig
	}{
		configuration.ExportConfig{
			Format: configuration.PRETTY,
		},
	}
	command := &cobra.Command{
		Use:     "deployment",
		Short:   "shows deployment data",
		Long:    "Shows deployment data",
		Example: "namespace deployment_names... [-n namespace_label]",
		Aliases: aliases,
		Run: func(command *cobra.Command, args []string) {
			deplData, err := func() (model.Renderer, error) {
				switch len(args) {
				case 0:
					logrus.Debugf("getting deployment from %q", ctx.Namespace)
					list, err := ctx.Client.GetDeploymentList(ctx.Namespace.ID)
					if err != nil {
						return nil, err
					}
					return list, nil
				default:
					deplNames := strset.NewSet(args)
					var showList deployment.DeploymentList = make([]deployment.Deployment, 0) // prevents panic
					list, err := ctx.Client.GetDeploymentList(ctx.Namespace.ID)
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
