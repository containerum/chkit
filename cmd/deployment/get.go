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
var GetDeployment = &cli.Command{
	Name:      "deployment",
	Aliases:   []string{"depl", "deployments", "deploy"},
	Usage:     "shows deployment data",
	ArgsUsage: "namespace [deployment_names ...]",
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
			deplNames := newSet(ctx.Args().Slice())
			var showList deployment.DeploymentList = make([]deployment.Deployment, 0) // prevents panic
			list, err := client.GetDeploymentList(namespace)
			if err != nil {
				return err
			}
			for _, depl := range list {
				if deplNames.Contain(depl.Name) {
					showList = append(showList, depl)
				}
			}
			show = showList
		}
		return util.WriteData(ctx, show)
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "file",
			Aliases: []string{"f"},
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
		},
	},
}

type set map[string]struct{}

func newSet(vals []string) set {
	var s set = make(map[string]struct{}, len(vals))
	for _, str := range vals {
		s[str] = struct{}{}
	}
	return s
}
func (s set) Contain(v string) bool {
	_, ok := s[v]
	return ok
}
