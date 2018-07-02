package clideployment

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/containerum/chkit/help"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/ninedraft/boxofstuff/str"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func DeleteVersion(ctx *context.Context) *cobra.Command {
	var flags struct {
		Deployment string `desc:"deployment name, can be chosen in interactive menu"`
		Version    string `desc:"deployment version, can be chosen in interactive menu"`
		Force      bool   `desc:"suppress confirmation" flag:"flag f"`
	}
	var command = &cobra.Command{
		Use:     "deployment-version",
		Aliases: []string{"depl-version", "devers", "deploy-vers", "depver", "deplver"},
		Short:   "delete inactive deployment version",
		Long:    help.GetString("delete deployment-version"),
		Example: "chkit delete deployment-version --deployment $DEPLOYMENT --version $VERSION --force",
		Run: func(cmd *cobra.Command, args []string) {
			var logger = ctx.Log.Command("run deployment version")
			var version semver.Version
			var deplName = str.Vector{flags.Deployment}.Append(args...).FirstNonEmpty()
			if deplName == "" {
				if flags.Force {
					ferr.Printf("deployment must be provided as arg or --deployment flag in force mode!\n")
					ctx.Exit(1)
				}
				logger.Debugf("getting deployment list")
				var deploymentList, err = ctx.Client.GetDeploymentList(ctx.GetNamespace().ID)
				if err != nil {
					ferr.Println(err)
					logger.WithError(err).Errorf("unable to get deployment list")
					ctx.Exit(1)
				}
				if deploymentList.Len() == 0 {
					ferr.Printf("You have no deployments in namespace %q!\n", ctx.GetNamespace())
					ctx.Exit(1)
				}
				logger.Debugf("selecting deployment")
				(&activekit.Menu{
					Title: "Select deployment",
					Items: activekit.ItemsFromIter(uint(deploymentList.Len()), func(index uint) *activekit.MenuItem {
						var depl = deploymentList[index]
						return &activekit.MenuItem{
							Label: depl.Name,
							Action: func() error {
								logger.Debugf("selected deployment %q", depl.Name)
								deplName = depl.Name
								return nil
							},
						}
					}),
				}).Run()
			}
			if flags.Version == "" {
				if flags.Force {
					ferr.Printf("version must be provided as --version flag in force mode!\n")
					ctx.Exit(1)
				}
				logger.Debugf("getting versions of deployment %q in namespace %q", deplName, ctx.GetNamespace())
				var deploymentVersions, err = ctx.Client.GetDeploymentVersions(ctx.GetNamespace().ID, deplName)
				if err != nil {
					ferr.Println(err)
					logger.WithError(err).Errorf("unable to get deployment versions")
					ctx.Exit(1)
				}
				deploymentVersions = deploymentVersions.Inactive()
				if deploymentVersions.Len() == 0 {
					ferr.Printf("Deployment %q in namespace %q has no inactive versions. Only inactive versions can be deleted!\n", deplName, ctx.GetNamespace())
					ctx.Exit(1)
				}
				logger.Debugf("selecting deployment version")
				(&activekit.Menu{
					Title: "Select version",
					Items: activekit.ItemsFromIter(uint(deploymentVersions.Len()), func(index uint) *activekit.MenuItem {
						var depl = deploymentVersions[index]
						return &activekit.MenuItem{
							Label: depl.Version.String(),
							Action: func() error {
								logger.Debugf("selected deployment version %q", depl.Version)
								version = depl.Version
								return nil
							},
						}
					}),
				}).Run()
			} else {
				var err error
				version, err = semver.ParseTolerant(flags.Version)
				if err != nil {
					ferr.Printf("invalid --version flag: %v\n", err)
					ctx.Exit(1)
				}
			}

			if flags.Force || activekit.YesNo("Are you sure you want to delete version %v of deployment %q?", version, deplName) {
				logger.Debugf("deleting version %q of deployment %q in namespace %q", version, deplName, ctx.GetNamespace())
				if err := ctx.Client.DeleteDeploymentVersion(ctx.GetNamespace().ID, deplName, version); err != nil {
					ferr.Println(err)
					logger.WithError(err).Errorf("unable to delete version %v of deployment %q", version, deplName)
					ctx.Exit(1)
				}
				fmt.Printf("Version %v of deployment %q in namespace %q deleted\n", version, deplName, ctx.GetNamespace())
			}
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
