package clideployment

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func DeleteContainer(ctx *context.Context) *cobra.Command {
	var flags struct {
		Force      bool   `desc:"suppress confirmation"`
		Deployment string `desc:"deployment name, required on --force"`
		Container  string `desc:"container name, required on --force"`
	}
	var command = &cobra.Command{
		Use:     "deployment-container",
		Aliases: []string{"depl-cont", "container", "dc"},
		Short:   "delete container",
		Long:    "Delete deployment container.",
		Run: func(cmd *cobra.Command, args []string) {
			var logger = ctx.Log.Command("delete deployment container")
			logger.Debugf("START")
			defer logger.Debugf("END")
			if flags.Force && flags.Deployment != "" {
				ferr.Printf("deployment name must be provided as --deployment while using --force")
				ctx.Exit(1)
			} else if flags.Deployment != "" {
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
						flags.Deployment = s
						return nil
					}),
				}).Run()
			}

			if flags.Container == "" && flags.Force {
				ferr.Printf("container name must be provided as --container while using --force")
				ctx.Exit(1)
			} else if flags.Container == "" {
				logger.Debugf("getting deployment %q", flags.Deployment)
				var depl, err = ctx.Client.GetDeployment(ctx.GetNamespace().ID, flags.Deployment)
				if err != nil {
					logger.WithError(err).Errorf("unable to get deployment %q", flags.Deployment)
					ferr.Println(err)
					ctx.Exit(1)
				}
				logger.Debugf("selecting container")
				(&activekit.Menu{
					Title: fmt.Sprintf("Select container in deployment %q", depl.Name),
					Items: activekit.StringSelector(depl.Containers.Names(), func(s string) error {
						flags.Container = s
						logger.Debugf("selected container %q", s)
						return nil
					}),
				}).Run()
			}
			if flags.Force || activekit.YesNo("Do you really want to delete container %q in deployment %q?", flags.Container, flags.Deployment) {
				logger.Debugf("deleting container %q in deployment %q", flags.Container, flags.Deployment)
				if err := ctx.Client.DeleteDeploymentContainer(ctx.GetNamespace().ID, flags.Deployment, flags.Container); err != nil {
					ferr.Println(err)
					logger.WithError(err).Debugf("unable to delete container %q in deployment %q", flags.Container, flags.Deployment)
					ctx.Exit(1)
				}
				fmt.Println("Ok")
			}
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
