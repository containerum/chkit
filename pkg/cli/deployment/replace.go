package clideployment

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"path/filepath"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/model/deployment/deplactive"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Replace(ctx *context.Context) *cobra.Command {
	var updFlags deplactive.UpdateFlags
	command := &cobra.Command{
		Use:     "deployment",
		Aliases: aliases,
		Short:   "replace deployment",
		Long: "Replaces deployment with new.\n" +
			"Has an one-line mode, suitable for integration with other tools,\n" +
			"and an interactive wizard mod",
		Run: func(cmd *cobra.Command, args []string) {
			var logger = ctx.Log.Command("update deployment")
			logger.Debugf("START")
			defer logger.Debugf("END")
			logger.Struct(updFlags)
			var depl deployment.Deployment
			var err error
			var deplName string
			switch len(args) {
			case 1:
				deplName = args[0]
			case 0:
				if updFlags.File != "" {
					depl, err = deplactive.FromFile(updFlags.File)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					updFlags = deplactive.FlagsFromDeployment(depl).UpdateFlags
				}
				if updFlags.Force {
					fmt.Printf("deployment name must be provided as first argument")
					os.Exit(1)
				}
				deplList, err := ctx.Client.GetDeploymentList(ctx.Namespace.ID)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				(&activekit.Menu{
					Title: "Select deployment",
					Items: activekit.StringSelector(deplList.Names(), func(s string) error {
						deplName = s
						return nil
					}),
				}).Run()
			}
			var flags = deplactive.Flags{
				Name:        deplName,
				UpdateFlags: updFlags,
			}
			depl, err = flags.Deployment()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			deplactive.Fill(&depl)
			if flags.Force {
				if err := deplactive.ValidateDeployment(depl); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if err := ctx.Client.ReplaceDeployment(ctx.Namespace.ID, depl); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Printf("Deployment %s created\n", depl.Name)
				return
			}
			logger.Debugf("getting configmap list")
			configmapList, err := ctx.Client.GetConfigmapList(ctx.Namespace.ID)
			if err != nil {
				logger.WithError(err).Errorf("unable to get configmap list")
				fmt.Println(err)
				os.Exit(1)
			}
			logger.Debugf("getting volume list")
			volumeList, err := ctx.Client.GetVolumeList(ctx.Namespace.ID)
			if err != nil {
				logger.WithError(err).Errorf("unable to get volume list")
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(depl.RenderTable())
			depl = deplactive.Wizard{
				EditName:   false,
				Deployment: &depl,
				Configmaps: configmapList.Names(),
				Volumes:    volumeList.Names(),
			}.Run()
			fmt.Println(depl.RenderTable())
			if !activekit.YesNo("Are you sure you want create deployment %s?", depl.Name) {
				(&activekit.Menu{
					Items: activekit.MenuItems{
						{
							Label: fmt.Sprintf("Save deployment %s to file", depl.Name),
							Action: func() error {
								for {
									var fname = activekit.Promt("Type filename (if ext is yaml or yml then file encodes as YAML, JSON by default): ")
									fname = strings.TrimSpace(fname)
									var err error
									var data string
									switch strings.ToLower(filepath.Ext(fname)) {
									case ".yaml", ".yml":
										fmt.Println("Encoding deployment as YAML")
										data, err = depl.RenderYAML()
									default:
										fmt.Println("Encoding deployment as JSON")
										data, err = depl.RenderJSON()
									}
									if err != nil {
										fmt.Println(err)
									}
									if err := ioutil.WriteFile(fname, []byte(data), os.ModePerm); err != nil {
										fmt.Println(err)
									}
								}
								return nil
							},
						},
					},
				}).Run()
				return
			}
			if err := ctx.Client.CreateDeployment(ctx.Namespace.ID, depl); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("Deployment %s created\n", depl.Name)
		},
	}
	if err := gpflag.ParseTo(&updFlags, command.Flags()); err != nil {
		angel.Angel(ctx, fmt.Errorf("it seems that the structure of the flags is set incorrectly: %v", err))
	}
	return command
}
