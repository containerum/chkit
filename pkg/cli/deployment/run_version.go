package clideployment

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/ninedraft/boxofstuff/str"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func RunVersion(ctx *context.Context) *cobra.Command {
	var flags struct {
		Deployment string `desc:"deployment name, can be chosen in interactive menu"`
		Version    string `desc:"deployment version, can be chosen in interactive menu.\nIf '-' or 'latest' then latest version is used."`
		Force      bool   `desc:"suppress confirmation" flag:"flag f"`
	}
	var command = &cobra.Command{
		Use:     "deployment-version",
		Aliases: []string{"depl-version", "devers", "deploy-vers", "depver", "deplver"},
		Short:   "run specific deployment version",
		// Long:    help.MustGetString("run deployment-version"),
		Example: "chkit run deployment-version --deployment $DEPLOYMENT --version $VERSION --force",
		Run: func(cmd *cobra.Command, args []string) {
			var logger = ctx.Log.Command("run deployment version")
			var version semver.Version
			var deplName = str.Vector{flags.Deployment}.Append(args...).FirstNonEmpty()
			if deplName == "" {
				if flags.Force {
					ferr.Printf("deployment must be provided as arg or --deployment flag in force mode!\n")
					ctx.Exit(1)
				}
				logger.Debugf("getting deployment list from namespace %q", ctx.GetNamespace())
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
				logger.Debugf("selecting deployment list")
				(&activekit.Menu{
					Title: "Select deployment",
					Items: activekit.ItemsFromIter(uint(deploymentList.Len()), func(index uint) *activekit.MenuItem {
						var depl = deploymentList[index]
						return &activekit.MenuItem{
							Label: depl.Name,
							Action: func() error {
								logger.Debugf("deployment %q selected", depl.Name)
								deplName = depl.Name
								return nil
							},
						}
					}),
				}).Run()
			}
			if flags.Version == "" || flags.Version == "-" || flags.Version == "latest" {
				if flags.Force {
					ferr.Printf("version must be provided as --version flag in force mode!\n")
					ctx.Exit(1)
				}
				logger.Debugf("getting deployment versions")
				var deploymentVersions, err = ctx.Client.GetDeploymentVersions(ctx.GetNamespace().ID, deplName)
				if err != nil {
					ferr.Println(err)
					logger.WithError(err).Errorf("unable to get deployment versions")
					ctx.Exit(1)
				}
				deploymentVersions = deploymentVersions.Inactive()
				if deploymentVersions.Len() == 0 {
					ferr.Printf("Deployment %q in namespace %q have no inactive versions\n", deplName, ctx.GetNamespace())
					ctx.Exit(1)
				}
				switch flags.Version {
				case "":
					logger.Debugf("selecting deployment version")
					(&activekit.Menu{
						Title: "Select version",
						Items: activekit.ItemsFromIter(uint(deploymentVersions.Len()), func(index uint) *activekit.MenuItem {
							var depl = deploymentVersions[index]
							return &activekit.MenuItem{
								Label: depl.Version.String(),
								Action: func() error {
									version = depl.Version
									return nil
								},
							}
						}),
					}).Run()
				case "-", "latest":
					version = deploymentVersions.SortByLess(func(a, b deployment.Deployment) bool {
						return a.Version.GT(b.Version)
					})[0].Version
				}
			} else {
				var err error
				version, err = semver.ParseTolerant(flags.Version)
				if err != nil {
					ferr.Printf("invalid --version flag: %v\n", err)
					ctx.Exit(1)
				}
			}

			if flags.Force || activekit.YesNo("Are you sure you want to run version %v of deployment %q?", version, deplName) {
				logger.Debugf("running version %q of deployment %q in namespace %q", version, deplName, ctx.GetNamespace())
				if err := ctx.Client.RunDeploymentVersion(ctx.GetNamespace().ID, deplName, version); err != nil {
					ferr.Println(err)
					logger.WithError(err).Errorf("unable to run version %v of deployment %q", version, deplName)
					ctx.Exit(1)
				}
				fmt.Printf("Version %v of deployment %q in namespace %q is started\n", version, deplName, ctx.GetNamespace())
			}

		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
