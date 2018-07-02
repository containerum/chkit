package clideployment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/containerum/chkit/help"
	"github.com/containerum/chkit/pkg/context"
	containerControl "github.com/containerum/chkit/pkg/controls/container"
	"github.com/containerum/chkit/pkg/export"
	"github.com/containerum/chkit/pkg/model/configmap"
	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/ninedraft/boxofstuff/str"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func CreateContainer(ctx *context.Context) *cobra.Command {
	var flags struct {
		Force      bool   `flag:"force f" desc:"suppress confirmation"`
		Name       string `desc:"container name, required on --force"`
		Deployment string `desc:"deployment name, required on --force"`
		File       string
		containerControl.Flags
	}
	command := &cobra.Command{
		Use:     "deployment-container",
		Aliases: []string{"depl-cont", "container", "dc"},
		Short:   "create deployment container.",
		Long:    help.GetString("create container"),
		Run: func(cmd *cobra.Command, args []string) {
			var logger = ctx.Log.Command("create deployment container")
			logger.Debugf("START")
			defer logger.Debugf("END")

			var cont container.Container

			if flags.File != "" {
				if flags.Deployment == "" && flags.Force {
					ferr.Printf("deployment name must be provided as --deployment while using --force")
					ctx.Exit(1)
				} else if flags.Deployment == "" {
					logger.Debugf("getting deployment list from namespace %q", ctx.GetNamespace())
					var depl, err = ctx.Client.GetDeploymentList(ctx.GetNamespace().ID)
					if err != nil {
						logger.WithError(err).Errorf("unable to get deployment list from namespace %q", ctx.GetNamespace())
						ferr.Println(err)
						ctx.Exit(1)
					}
					logger.Debugf("selecting deployment")
					(&activekit.Menu{
						Title: "Select deployment",
						Items: activekit.StringSelector(depl.Names(), func(s string) error {
							logger.Debugf("deployment %q selected", s)
							flags.Deployment = s
							return nil
						}),
					}).Run()
				}
				if flags.Name == "" && flags.Force {
					if flags.Image == "" {
						ferr.Printf("container --image must be provided while using --force")
						ctx.Exit(1)
					}
					flags.Name = str.Vector{namegen.Color(), container.ImageName(flags.Image)}.Join("-")
				}
				var err error
				logger.Debugf("building container from flags")
				cont, err = flags.Container()
				if err != nil {
					logger.WithError(err).Errorf("unable to build container from flags")
					ferr.Println(err)
					ctx.Exit(1)
				}
				cont.Name = flags.Name
			} else {
				var data, err = ioutil.ReadFile(flags.File)
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				switch filepath.Ext(flags.File) {
				case ".yaml", ".yml":
					err = yaml.Unmarshal(data, &cont)
				default:
					err = json.Unmarshal(data, &cont)
				}
				if err != nil {
					ferr.Printf("Unable to import container from file %q:\n%v\n", flags.File, err)
					ctx.Exit(1)
				}
			}
			cont = containerControl.Default(cont)
			fmt.Println(cont.RenderTable())
			if flags.Force {
				logger.Debugf("running in --force mode")
				cont = containerControl.Default(cont)
				logger.Debugf("validating changed container %q", cont.Name)
				if err := cont.Validate(); err != nil {
					logger.WithError(err).Errorf("invalid container %q", cont.Name)
					ferr.Println(err)
					ctx.Exit(1)
				}
				logger.Debugf("creating container %q", cont.Name)
				if err := ctx.Client.CreateDeploymentContainer(ctx.GetNamespace().ID, flags.Deployment, cont); err != nil {
					logger.WithError(err).Errorf("unable to replace container %q", cont.Name)
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Println("Ok")
				return
			}

			//	volumes, err := ctx.Client.GetVolumeList(ctx.GetNamespace().ID)
			//	if err != nil {
			//		ferr.Println(err)
			//		os.Exit(1)
			//	}

			var deployments = make(chan deployment.DeploymentList)
			go func() {
				logger := logger.Component("getting deployment list")
				logger.Debugf("START")
				defer logger.Debugf("END")
				defer close(deployments)
				deplList, err := ctx.Client.GetDeploymentList(ctx.GetNamespace().ID)
				if err != nil {
					logger.WithError(err).Errorf("unable to get deployment list from namespace %q", ctx.GetNamespace())
					ferr.Println(err)
					ctx.Exit(1)
				}
				deployments <- deplList
			}()

			var configs = make(chan configmap.ConfigMapList)
			go func() {
				logger := logger.Component("getting configmap list")
				logger.Debugf("START")
				defer logger.Debugf("END")
				defer close(configs)
				configList, err := ctx.Client.GetConfigmapList(ctx.GetNamespace().ID)
				if err != nil {
					logger.WithError(err).Errorf("unable to get configmap list")
					ferr.Println(err)
					ctx.Exit(1)
				}
				configs <- configList
			}()
			logger.Debugf("running wizard")
			cont = containerControl.Wizard{
				EditName:   true,
				Container:  cont,
				Deployment: flags.Deployment,
				//		Volumes:     volumes.Names(),
				Configs:     (<-configs).Names(),
				Deployments: (<-deployments).Names(),
			}.Run()
			fmt.Println(cont.RenderTable())
			if activekit.YesNo("Are you sure you want to create container %q in deployment %q?", cont.Name, flags.Deployment) {
				logger.Debugf("creating container %q in deployment %q", cont.Name, flags.Deployment)
				if err := ctx.Client.CreateDeploymentContainer(ctx.GetNamespace().ID, flags.Deployment, cont); err != nil {
					logger.WithError(err).Errorf("unable to replace container %q in deployment %q", cont.Name, flags.Deployment)
					ferr.Println(err)
				}
				fmt.Println("Ok")
			}
			(&activekit.Menu{
				Items: activekit.MenuItems{
					{
						Label: "Save container to file",
						Action: func() error {
							for {
								var fname = activekit.Promt("Type output filename: ")
								export.ExportData(cont, export.ExportConfig{
									Filename: fname,
									Format:   export.YAML,
								})
							}

							return nil
						},
					},
				},
			}).Run()
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
