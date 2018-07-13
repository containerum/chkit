package clideployment

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
	containerControl "github.com/containerum/chkit/pkg/controls/container"
	"github.com/containerum/chkit/pkg/model/configmap"
	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/porta"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func ReplaceContainer(ctx *context.Context) *cobra.Command {
	var flags struct {
		Force         bool   `flag:"force f" desc:"suppress confirmation"`
		ContainerName string `flag:"container" desc:"container name, required on --force"`
		Deployment    string `desc:"deployment name, required on --force"`
		containerControl.ReplaceFlags
		porta.Importer
	}
	command := &cobra.Command{
		Use:     "deployment-container",
		Aliases: []string{"depl-cont", "container", "dc"},
		Short:   "Replace deployment container.",
		Long: "Replace deployment container.\n" +
			"Runs in one-line mode, suitable for integration with other tools, and in interactive wizard mode.",
		Run: func(cmd *cobra.Command, args []string) {
			var logger = ctx.Log.Command("replace deployment container")
			logger.Debugf("START")
			defer logger.Debugf("END")
			if flags.Force {
				if flags.Deployment == "" {
					ferr.Printf("deployment name must be provided as --deployment while using --force")
					ctx.Exit(1)
				}
				if flags.ContainerName == "" {
					ferr.Printf("container name must be provided as --container while using --force")
					ctx.Exit(1)
				}
				var depl, err = ctx.Client.GetDeployment(ctx.GetNamespace().ID, flags.Deployment)
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				cont, ok := depl.Containers.GetByName(flags.ContainerName)
				if !ok {
					ferr.Printf("container %q doesn't exist", flags.ContainerName)
					ctx.Exit(1)
				}
				if flags.ImportActivated() {
					var importedCont container.Container
					if err := flags.Import(&importedCont); err != nil {
						ferr.Println(err)
						ctx.Exit(1)
					}
					cont, err = flags.Patch(importedCont)
				} else {
					flagCont, err := flags.Container()
					if err != nil {
						ferr.Println(err)
						ctx.Exit(1)
					}
					cont, err = flags.Patch(flagCont)
				}
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				cont.Name = flags.ContainerName
				depl.Containers, _ = depl.Containers.Replace(cont)
				if err := ctx.Client.ReplaceDeployment(ctx.GetNamespace().ID, depl); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Println("Ok")
				return
			}

			var depl deployment.Deployment
			if flags.Deployment == "" {
				logger.Debugf("getting deployment list from namespace %q", ctx.GetNamespace())
				deplList, err := ctx.Client.GetDeploymentList(ctx.GetNamespace().ID)
				if err != nil {
					logger.WithError(err).Errorf("unable to get deployment list from namespace %q", ctx.GetNamespace())
					ferr.Println(err)
					ctx.Exit(1)
				}
				logger.Debugf("selecting deployment")
				(&activekit.Menu{
					Title: "Select deployment",
					Items: activekit.ItemsFromIter(uint(deplList.Len()), func(index uint) *activekit.MenuItem {
						var d = deplList[index]
						return &activekit.MenuItem{
							Label: d.Name,
							Action: func() error {
								flags.Deployment = d.Name
								depl = d
								logger.Debugf("deployment %q selected", d.Name)
								return nil
							},
						}
					}),
				}).Run()
			} else {
				logger.Debugf("getting deployment %q", flags.Deployment)
				var err error
				depl, err = ctx.Client.GetDeployment(ctx.GetNamespace().ID, flags.Deployment)
				if err != nil {
					logger.WithError(err).Errorf("unable to get deployment %q", flags.Deployment)
					ferr.Println(err)
					ctx.Exit(1)
				}
			}

			var cont container.Container

			if flags.ContainerName == "" {
				logger.Debugf("selecting container")
				(&activekit.Menu{
					Title: fmt.Sprintf("Select container in deployment %q", depl.Name),
					Items: activekit.ItemsFromIter(uint(len(depl.Containers)), func(index uint) *activekit.MenuItem {
						var c = depl.Containers[index]
						return &activekit.MenuItem{
							Label: c.Name,
							Action: func() error {
								flags.ContainerName = c.Name
								logger.Debugf("selected container %q", c.Name)
								cont = c
								return nil
							},
						}
					}),
				}).Run()
			} else {
				ok := false
				cont, ok = depl.Containers.GetByName(flags.ContainerName)
				if !ok {
					ferr.Printf("container %q not found in deployment %q", flags.ContainerName, depl.Name)
					ctx.Exit(1)
				}
			}

			if flags.ImportActivated() {
				var importedCont container.Container
				if err := flags.Import(&importedCont); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				cont = cont.Patch(importedCont)
			} else {
				flagCont, err := flags.Container()
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				cont = cont.Patch(flagCont)
			}

			logger.Debugf("building container from flags")
			cont, err := flags.Patch(cont)
			if err != nil {
				logger.WithError(err).Errorf("unable to build container from flags")
				ferr.Println(err)
				ctx.Exit(1)
			}
			cont.Name = flags.ContainerName
			fmt.Println(cont.RenderTable())

			/*var volumes = make(chan volume.VolumeList)
			go func() {
				logger := logger.Component("getting namespace list")
				logger.Debugf("START")
				defer logger.Debugf("END")
				defer close(volumes)
				var volumeList, err = ctx.Client.GetVolumeList(ctx.GetNamespace().ID)
				if err != nil {
					logger.WithError(err).Errorf("unable to get volume list from namespace %q", ctx.GetNamespace())
					ferr.Println(err)
					ctx.Exit(1)
				}
				volumes <- volumeList
			}()*/

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
				Container:  cont,
				Deployment: flags.Deployment,
				//	Volumes:     (<-volumes).Names(),
				Configs:     (<-configs).Names(),
				Deployments: (<-deployments).Names(),
			}.Run()

			if activekit.YesNo("Are you sure you want to update container %q in deployment %q?", cont.Name, flags.Deployment) {
				logger.Debugf("replacing container %q in deployment %q", cont.Name, flags.Deployment)
				if err := ctx.Client.ReplaceDeploymentContainer(ctx.GetNamespace().ID, flags.Deployment, cont); err != nil {
					logger.WithError(err).Errorf("unable to replace container %q in deployment %q", cont.Name, flags.Deployment)
					ferr.Println(err)
				}
			}
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
