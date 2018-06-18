package clideployment

import (
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	deployment2 "github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func GetVersions(ctx *context.Context) *cobra.Command {
	var flags struct {
		LastN   uint64                     `desc:"limit n versions to show"`
		Output  configuration.ExportFormat `desc:"output format, json/yaml"`
		File    string                     `desc:"output file, optional, default is STDOUT"`
		Version string                     `desc:"version query, examples: <1.0.0, <=1.0.0, !1.0.0"`
	}
	var command = &cobra.Command{
		Use:     "deployment-versions",
		Aliases: []string{"depl-ver", "depvers", "deployment-version"},
		Short:   "get deployment versions",
		Example: "chkit get deployment-versions MY_DEPLOYMENT [--last-n 4] [--version >=1.0.0] [--output yaml] [--file versions.yaml]",
		Long: "Get deployment versions.\n" +
			"You can filter versions by specifying version query (--version):\n" +
			// this part of docs is adapted comments from github.com/blang/semver
			"Valid queries are:\n\n" +
			"  - \"<1.0.0\"\n" +
			"  - \"<=1.0.0\"\n" +
			"  - \">1.0.0\"\n" +
			"  - \">=1.0.0\"\n" +
			"  - \"1.0.0\", \"=1.0.0\", \"==1.0.0\"\n" +
			"  - \"!1.0.0\", \"!=1.0.0\"\n" +
			"A query can consist of multiple querys separated by space:\n" +
			"queries can be linked by logical AND:\n" +
			"  - \">1.0.0 <2.0.0\" would match between both querys, so \"1.1.1\" and \"1.8.7\" but not \"1.0.0\" or \"2.0.0\"\n" +
			"  - \">1.0.0 <3.0.0 !2.0.3-beta.2\" would match every version between 1.0.0 and 3.0.0 except 2.0.3-beta.2\n" +
			"Queries can also be linked by logical OR:\n" +
			"  - \"<2.0.0 || >=3.0.0\" would match \"1.x.x\" and \"3.x.x\" but not \"2.x.x\"\n" +
			"AND has a higher precedence than OR. It's not possible to use brackets.\n" +
			"Queries can be combined by both AND and OR\n" +
			" - `>1.0.0 <2.0.0 || >3.0.0 !4.2.1` would match `1.2.3`, `1.9.9`, `3.1.1`, but not `4.2.1`, `2.1.1`",
		Run: func(cmd *cobra.Command, args []string) {
			var logger = ctx.Log.Command("get deployment-versions")
			logger.Debugf("START")
			defer logger.Debugf("END")
			logger.StructFields(flags)
			var deployment string
			switch len(args) {
			case 0:
				logger.Debugf("getting deployment list")
				var list, err = ctx.Client.GetDeploymentList(ctx.Namespace.ID)
				if err != nil {
					logger.WithError(err).Debugf("unable to get deployment list")
					fmt.Println(err)
					os.Exit(1)
				}
				logger.Debugf("selecting deployment")
				(&activekit.Menu{
					Title: "Select deployment",
					Items: activekit.StringSelector(list.Names(), func(s string) error {
						logger.Debugf("using deployment %q", s)
						deployment = s
						return nil
					}),
				}).Run()
			case 1:
				deployment = args[0]
				logger.Debugf("using deployment %q", args[0])
			default:
				cmd.Help()
				os.Exit(1)
			}
			logger.Debugf("getting versions of deployment %q", deployment)
			var versions, err = ctx.Client.GetDeploymentVersions(ctx.Namespace.ID, deployment)
			if err != nil {
				logger.WithError(err).Errorf("unable to get versions of deployment %q", deployment)
				fmt.Println(err)
				os.Exit(1)
			}
			logger.Debugf("retrieved %d versions", len(versions))
			if flags.Version != "" {
				logger.Debugf("parsing versions query %q", flags.Version)
				query, err := semver.ParseRange(flags.Version)
				if err != nil {
					logger.WithError(err).Errorf("unable to parse version query")
					fmt.Println(err)
					os.Exit(1)
				}
				logger.Debugf("selecting deployments by query %q", flags.Version)
				versions = versions.Filter(func(depl deployment2.Deployment) bool {
					return query(depl.Version)
				})
				logger.Debugf("%d versions selected", len(versions))
			}
			if flags.LastN > 0 && flags.LastN < uint64(len(versions)) {
				logger.Debugf("selecting last %d versions", flags.LastN)
				versions = versions[:flags.LastN]
			}
			logger.Debugf("exporting versions data")
			if err := configuration.ExportData(versions, configuration.ExportConfig{
				Filename: flags.File,
				Format:   flags.Output,
			}); err != nil {
				logger.WithError(err).Errorf("unable to export versions data")
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
