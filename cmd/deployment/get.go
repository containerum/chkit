package clideployment

import (
	"strings"

	"github.com/containerum/chkit/cmd/cmdutil"
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
	Description: "shows deployment data. Aliases: " + strings.Join(aliases, ", "),
	UsageText:   "namespace deployment_names... [-n namespace_label]",
	Action: func(ctx *cli.Context) error {
		if ctx.Bool("help") {
			return cli.ShowSubcommandHelp(ctx)
		}
		client := cmdutil.GetClient(ctx)
		defer cmdutil.StoreClient(ctx, client)

		var show model.Renderer
		switch ctx.NArg() {
		case 0:
			namespace := cmdutil.GetNamespace(ctx)
			logrus.Debugf("getting deployment from %q", namespace)
			list, err := client.GetDeploymentList(namespace)
			if err != nil {
				return err
			}
			show = list
		default:
			namespace := cmdutil.GetNamespace(ctx)
			deplNames := cmdutil.NewSet(ctx.Args().Slice())
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
		return cmdutil.ExportDataCommand(ctx, show)
	},
	Flags: cmdutil.GetFlags,
}
