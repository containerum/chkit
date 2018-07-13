package diff

import (
	"io/ioutil"
	"os"

	"github.com/blang/semver"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/containerum/chkit/pkg/util/prettydiff"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Get(ctx *context.Context) *cobra.Command {
	var flags struct {
		Deployment     string `desc:"deployment name, optional"`
		Version        string `desc:"first deployment version to compare"`
		AnotherVersion string `desc:"second deployment version to compare"`
		Output         string `desc:"diff output, STDOUT by default"`
	}
	var command = &cobra.Command{
		Use:   "diff",
		Short: "show diff between deployment versions",
		Run: func(cmd *cobra.Command, args []string) {
			if flags.Deployment == "" {
				var deplList, err = ctx.Client.GetDeploymentList(ctx.GetNamespace().ID)
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				if deplList.Len() == 0 {
					ferr.Printf("You have no deployments in namespace %q!", ctx.GetNamespace())
				}
				var depl = selectDeploment(ctx, deplList)
				flags.Deployment = depl.Name
			}
			var version, anotherVersion semver.Version
			var versions []semver.Version

			if flags.Version == "" || flags.Version == "-" {
				if versions == nil {
					var deploymentVersions, err = ctx.Client.GetDeploymentVersions(ctx.GetNamespace().ID, flags.Deployment)
					if err != nil {
						ferr.Println(err)
						ctx.Exit(1)
					}
					versions = deploymentVersions.Versions()
				}
				if len(versions) == 0 {
					ferr.Printf("you have no additional versions!\n")
					ctx.Exit(1)
				}
				switch flags.Version {
				case "-":
					version = versions[0]
				default:
					version = selectVersion(ctx, "Select version", versions)
				}

			}

			var diff string

			if flags.AnotherVersion != "" {
				var v, err = semver.ParseTolerant(flags.AnotherVersion)
				if err != nil {
					ferr.Printf("unable to parse --another-version: %v", err)
				}
				anotherVersion = v
				diff, err = ctx.Client.GetDeploymentDiffBetweenVersions(ctx.GetNamespace().ID, flags.Deployment, version, anotherVersion)
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
			} else {
				var err error
				diff, err = ctx.Client.GetDeploymentDiffWithPreviousVersion(ctx.GetNamespace().ID, flags.Deployment, version)
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
			}

			switch flags.Output {
			case "", "-":
				prettydiff.Fprint(os.Stdout, diff)
			default:
				if err := ioutil.WriteFile(flags.Output, []byte(diff), os.ModePerm); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
			}
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}

func selectDeploment(ctx *context.Context, list deployment.DeploymentList) deployment.Deployment {
	var selectedDepl deployment.Deployment
	(&activekit.Menu{
		Title: "Select deployment",
		Items: activekit.ItemsFromIter(uint(list.Len()),
			func(index uint) *activekit.MenuItem {
				var depl = list[index]
				return &activekit.MenuItem{
					Label: depl.Name,
					Action: func() error {
						selectedDepl = depl
						return nil
					},
				}
			}),
	}).Run()
	return selectedDepl.Copy()
}

func selectVersion(ctx *context.Context, title string, versions []semver.Version) semver.Version {
	var selectedVersion semver.Version
	(&activekit.Menu{
		Title: title,
		Items: activekit.ItemsFromIter(uint(len(versions)), func(index uint) *activekit.MenuItem {
			var version = versions[index]
			return &activekit.MenuItem{
				Label: version.String(),
				Action: func() error {
					selectedVersion = version
					return nil
				},
			}
		}),
	}).Run()
	return selectedVersion
}
