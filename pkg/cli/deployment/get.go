package clideployment

import (
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/porta"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/ninedraft/boxofstuff/strset"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	ErrNoNamespaceSpecified chkitErrors.Err = "no namespace specified"
)

var aliases = []string{"depl", "deployments", "deploy", "de", "dpl", "depls", "dep"}

func Get(ctx *context.Context) *cobra.Command {
	var flags struct {
		Solution string
		porta.Exporter
	}
	command := &cobra.Command{
		Use:     "deployment",
		Short:   "show deployment data",
		Long:    "Print deployment data.",
		Example: "get deployment_names... [-n project_label]",
		Aliases: aliases,
		Run: func(command *cobra.Command, args []string) {
			var logger = ctx.Log.Command("get deployment")
			logger.Debugf("START")
			defer logrus.Debugf("END")
			var deplData model.Renderer

			switch len(args) {
			case 1:
				if flags.Solution != "" {
					logger.Debugf("getting deployment list from namespace %q and solution %q",
						ctx.GetNamespace(), flags.Solution)
					var depls, err = ctx.Client.GetSolutionDeployments(ctx.GetNamespace().ID, flags.Solution)
					if err != nil {
						logger.WithError(err).Errorf("unable to get deployment list from namespace %q and solution %q", ctx.GetNamespace(), flags.Solution)
						ferr.Println(err)
						ctx.Exit(1)
					}
					var depl, ok = depls.GetByName(args[0])
					if !ok {
						ferr.Printf("deployment %q not found\n", args[0])
					}
					deplData = depl
				} else {
					logger.Debugf("getting deployment list from namespace %q", ctx.GetNamespace())
					var depl, err = ctx.Client.GetDeployment(ctx.GetNamespace().ID, args[0])
					if err != nil {
						logger.WithError(err).Errorf("unable to get deployment list from namespace %q", ctx.GetNamespace())
						ferr.Println(err)
						ctx.Exit(1)
					}
					deplData = depl
				}
			default:
				var deplList deployment.DeploymentList
				if flags.Solution != "" {
					var depls, err = ctx.Client.GetSolutionDeployments(ctx.GetNamespace().ID, flags.Solution)
					if err != nil {
						logger.WithError(err).Errorf("unable to get deployment list from namespace %q and solution %q", ctx.GetNamespace(), flags.Solution)
						ferr.Println(err)
						ctx.Exit(1)
					}
					deplList = depls
				} else {
					var depls, err = ctx.Client.GetDeploymentList(ctx.GetNamespace().ID)
					if err != nil {
						logger.WithError(err).Errorf("unable to get deployment list from namespace %q", ctx.GetNamespace())
						ferr.Println(err)
						ctx.Exit(1)
					}
					deplList = depls
				}
				if len(args) > 0 {
					var deplNames = strset.NewSet(args)
					deplList = deplList.Filter(func(depl deployment.Deployment) bool {
						return deplNames.Have(depl.Name)
					})
				}
				deplData = deplList
			}
			logger.Debugf("exporting deployment data")
			if err := flags.Export(deplData); err != nil {
				ferr.Println(err)
				ctx.Exit(1)
			}
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
