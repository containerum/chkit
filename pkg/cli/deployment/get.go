package clideployment

import (
	"os"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
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
				var depl, err = ctx.Client.GetDeployment(ctx.Namespace.ID, args[0])
				if err != nil {
					ferr.Println(err)
					os.Exit(1)
				}
				deplData = depl
			} else {
				var list, err = ctx.Client.GetDeploymentList(ctx.Namespace.ID)
				if err != nil {
					ferr.Println(err)
					os.Exit(1)
				}
				var names = strset.NewSet(args)
				list = list.Filter(func(depl deployment.Deployment) bool {
					return names.Have(depl.Name)
				})
				deplData = list
			}
			if err := configuration.ExportData(deplData, configuration.ExportConfig{
				Format:   configuration.ExportFormat(flags.Output),
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
	return command
}
