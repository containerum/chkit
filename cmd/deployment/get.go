package clideployment

import (
	"fmt"
	"strings"
	"time"

	"github.com/containerum/chkit/cmd/cmdutil"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/animation"
	"github.com/containerum/chkit/pkg/util/trasher"
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

		anime := &animation.Animation{
			Framerate:      0.5,
			Source:         trasher.NewSilly(),
			ClearLastFrame: true,
		}
		go func() {
			time.Sleep(4 * time.Second)
			anime.Run()
		}()
		if rend, err := getDeployment(ctx, anime); err != nil {
			return err
		} else {
			return cmdutil.ExportDataCommand(ctx, rend)
		}
	},
	Flags: cmdutil.GetFlags,
}

func getDeployment(ctx *cli.Context, anime *animation.Animation) (model.Renderer, error) {
	defer anime.Stop()
	client := cmdutil.GetClient(ctx)
	switch ctx.NArg() {
	case 0:
		namespace := cmdutil.GetNamespace(ctx)
		logrus.Debugf("getting deployment from %q", namespace)
		list, err := client.GetDeploymentList(namespace)
		if err != nil {
			logrus.WithError(err).Errorf("unable to get list of deployments from %q", namespace)
			fmt.Printf("Unable to get list of deployments from namespace %q\n", namespace)
			return nil, err
		}
		return list, nil
	case 1:
		namespace := cmdutil.GetNamespace(ctx)
		deplName := ctx.Args().First()
		depl, err := client.GetDeployment(namespace, deplName)
		if err != nil {
			logrus.WithError(err).Errorf("unable to get deployment %q from namespace %q", deplName, namespace)
			fmt.Printf("Unable to get deployment %q from namespace %q\n", deplName, namespace)
			return nil, err
		}
		return depl, nil
	default:
		namespace := cmdutil.GetNamespace(ctx)
		deplNames := cmdutil.NewSet(ctx.Args().Slice())
		var showList deployment.DeploymentList = make([]deployment.Deployment, 0) // prevents panic
		list, err := client.GetDeploymentList(namespace)
		if err != nil {
			logrus.WithError(err).Errorf("unable to get list of deployments")
			if len(deplNames) == 1 {
				logrus.WithError(err).Errorf("unable to get list of deployments from %q", namespace)
			} else {
				fmt.Printf("Unable to get deployments %v from namespace %q", deplNames.Slice(), namespace)
			}
			return nil, err
		}
		for _, depl := range list {
			if deplNames.Have(depl.Name) {
				showList = append(showList, depl)
			}
		}
		return list, nil
	}
}
