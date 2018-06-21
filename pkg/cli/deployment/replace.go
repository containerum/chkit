package clideployment

import (
	"fmt"

	"os"

	"github.com/containerum/chkit/pkg/context"
	containerControl "github.com/containerum/chkit/pkg/controls/container"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Replace(ctx *context.Context) *cobra.Command {
	var flags struct {
		Force         bool   `flag:"force f" desc:"suppress confirmation"`
		ContainerName string `flag:"container" desc:"container name"`
		Deployment    string `desc:"deployment name"`
		containerControl.Flags
	}
	command := &cobra.Command{
		Use:     "deployment-container",
		Aliases: aliases,
		Short:   "Replace deployment container.",
		Long: "Replace deployment container.\n" +
			"Runs in one-line mode, suitable for integration with other tools, and in interactive wizard mode.",
		Run: func(cmd *cobra.Command, args []string) {
			var logger = ctx.Log.Command("replace deployment container")
			logger.Debugf("START")
			defer logger.Debugf("END")

			if flags.Deployment == "" && flags.Force {
				ferr.Printf("deployment name must be provided as --deployment while using --force")
				os.Exit(1)
			} else if flags.Deployment == "" {
				var depl, err = ctx.Client.GetDeploymentList(ctx.Namespace.ID)
				if err != nil {
					ferr.Println(err)
					os.Exit(1)
				}
				(&activekit.Menu{
					Title: "Select deployment",
					Items: activekit.StringSelector(depl.Names(), func(s string) error {
						flags.Deployment = s
						return nil
					}),
				}).Run()
			}

			if flags.ContainerName == "" && flags.Force {
				ferr.Printf("container name must be provided as --container while using --force")
				os.Exit(1)
			} else if flags.ContainerName == "" {
				var depl, err = ctx.Client.GetDeployment(ctx.Namespace.ID, flags.Deployment)
				if err != nil {
					ferr.Println(err)
					os.Exit(1)
				}
				(&activekit.Menu{
					Title: fmt.Sprintf("Select container in deployment %q", depl.Name),
					Items: activekit.StringSelector(depl.Containers.Names(), func(s string) error {
						flags.ContainerName = s
						return nil
					}),
				}).Run()
			}
			var cont, err = flags.Container()
			if err != nil {
				ferr.Println(err)
				os.Exit(1)
			}
			cont.Name = flags.ContainerName
			fmt.Println(cont)

			//	volumes, err := ctx.Client.GetVolumeList(ctx.Namespace.ID)
			//	if err != nil {
			//		ferr.Println(err)
			//		os.Exit(1)
			//	}

			deployments, err := ctx.Client.GetDeploymentList(ctx.Namespace.ID)
			if err != nil {
				ferr.Println(err)
				os.Exit(1)
			}

			configs, err := ctx.Client.GetConfigmapList(ctx.Namespace.ID)
			if err != nil {
				ferr.Println(err)
				os.Exit(1)
			}
			containerControl.Wizard{
				Container:  cont,
				Deployment: flags.Deployment,
				//		Volumes:     volumes.Names(),
				Configs:     configs.Names(),
				Deployments: deployments.Names(),
			}.Run()
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
