package clideployment

import (
	"os"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/export"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/ninedraft/boxofstuff/strset"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	ErrNoNamespaceSpecified chkitErrors.Err = "no namespace specified"
)

var aliases = []string{"depl", "deployments", "deploy"}

func Get(ctx *context.Context) *cobra.Command {
	var flags struct {
		File   string `desc:"output file, STDOUT by default"`
		Output string `flag:"output o" desc:"output format, json/yaml"`
	}
	command := &cobra.Command{
		Use:     "deployment",
		Short:   "show deployment data",
		Long:    "Print deployment data.",
		Example: "namespace deployment_names... [-n namespace_label]",
		Aliases: aliases,
		Run: func(command *cobra.Command, args []string) {
			var logger = ctx.Log.Command("get deployment")
			logger.Debugf("START")
			defer logrus.Debugf("END")
			var deplData model.Renderer
			if len(args) == 1 {
				logger.Debugf("getting deployment %q from namespace %q", args[0], ctx.Namespace)
				var depl, err = ctx.Client.GetDeployment(ctx.Namespace.ID, args[0])
				if err != nil {
					logger.WithError(err).Errorf("unable to get deployment %q from namespace %q", args[0], ctx.Namespace)
					ferr.Println(err)
					os.Exit(1)
				}
				deplData = depl
			} else {
				logrus.Debugf("getting deployment list from namespace %q", ctx.Namespace)
				var list deployment.DeploymentList
				var err error
				if solutionName, _ := command.Flags().GetString("solution_name"); solutionName != "" {
					list, err = ctx.Client.GetSolutionDeployments(ctx.Namespace.ID, solutionName)
				} else {
					list, err = ctx.Client.GetDeploymentList(ctx.Namespace.ID)
				}
				if err != nil {
					logger.WithError(err).Errorf("unable to get deployment list from namespace %q", ctx.Namespace)
					ferr.Println(err)
					os.Exit(1)
				}
				if len(args) > 0 {
					logrus.Debugf("filtering deployment list: including only %v", args)
					var names = strset.NewSet(args)
					list = list.Filter(func(depl deployment.Deployment) bool {
						return names.Have(depl.Name)
					})
				}
				deplData = list
			}
			logger.Debugf("exporting deployment data")
			if err := export.ExportData(deplData, export.ExportConfig{
				Format:   export.ExportFormat(flags.Output),
				Filename: flags.File,
			}); err != nil {
				logrus.WithError(err).Errorf("unable to export data")
				angel.Angel(ctx, err)
			}
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	command.PersistentFlags().
		StringP("solution_name", "s", "", "solution name")
	return command
}
