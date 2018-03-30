package clideployment

import (
	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

var (
	ErrNoNamespaceSpecified chkitErrors.Err = "no namespace specified"
)

var aliases = []string{"depl", "deployments", "deploy"}

var GetDeployment = &cli.Command{
	Name:        "deployment",
	Aliases:     aliases,
	Usage:       "shows deployment data",
	Description: "shows deployment data. Aliases: depl, deployments, deploy",
	ArgsUsage:   "namespace [deployment_names ...]",
	Action: func(ctx *cli.Context) error {
		if ctx.Bool("help") {
			return cli.ShowSubcommandHelp(ctx)
		}
		client := util.GetClient(ctx)
		defer util.StoreClient(ctx, client)

		var show model.Renderer
		switch ctx.NArg() {
		case 0:
			namespace := util.GetNamespace(ctx)
			logrus.Debugf("getting deployment from %q", namespace)
			list, err := client.GetDeploymentList(namespace)
			if err != nil {
				return err
			}
			show = list
		default:
			namespace := util.GetNamespace(ctx)
			deplNames := util.NewSet(ctx.Args().Slice())
			var showList deployment.DeploymentList = make([]deployment.Deployment, 0) // prevents panic
			list, err := client.GetDeploymentList(namespace)
			if err != nil {
				return err
			}
			for _, depl := range list {
				if deplNames.Have(depl.Name) {
					showList = append(showList, depl)
				}
			}
			show = showList
		}
		return util.WriteData(ctx, show)
	},
	Flags: util.GetFlags,
}
