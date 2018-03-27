package clideployment

import (
	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/deployment"
	"gopkg.in/urfave/cli.v2"
)

var (
	ErrNoNamespaceSpecified chkitErrors.Err = "no namespace specified"
)
var GetDeployment = &cli.Command{
	Name:      "deployment",
	Usage:     "shows deployment data",
	ArgsUsage: "namespace [deployment_names ...]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "json",
			Usage: "writes json data to file. If filename is \"stdout\" then prints straight to std output",
		},
		&cli.StringFlag{
			Name:  "yaml",
			Usage: "writes yaml data to file. If filename is \"stdout\" then prints straight to std output",
		},
	},
	Action: func(ctx *cli.Context) error {
		if ctx.Bool("help") {
			return cli.ShowSubcommandHelp(ctx)
		}
		client := util.GetClient(ctx)
		log := util.GetLog(ctx)

		defer func() {
			log := util.GetLog(ctx)
			util.SetClient(ctx, client)
			log.Debugf("writing tokens to disk")
			err := util.SaveTokens(ctx, client.Tokens)
			if err != nil {
				log.Debugf("error while saving tokens: %v", err)
				panic(err)
			}
		}()

		var show model.Renderer
		switch ctx.NArg() {
		case 0:
			cli.ShowCommandHelpAndExit(ctx, "deployment", 2)
		case 1:
			namespace := ctx.Args().First()
			log.Debugf("getting deployment from %q", namespace)
			list, err := client.GetDeploymentList(namespace)
			if err != nil {
				return err
			}
			show = list
		case 2:
			namespace := ctx.Args().First()
			deplName := ctx.Args().Slice()[1]
			depl, err := client.GetDeployment(namespace, deplName)
			if err != nil {
				return err
			}
			show = depl
		default:
			namespace := ctx.Args().First()
			deplNames := newSet(ctx.Args().Tail())
			var showList deployment.DeploymentList = make([]deployment.Deployment, len(deplNames))
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
