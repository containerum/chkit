package clitemplate

import (
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/export"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var aliases_envs = []string{"template_env", "tmpl_env", "envs", "environments", "templates_environments", "tmpls_env", "tmpenv", "tmp_env", "tmps_env", "tmpsenv"}

func GetEnvs(ctx *context.Context) *cobra.Command {
	exportConfig := export.ExportConfig{}
	command := &cobra.Command{
		Use:     "template_envs",
		Aliases: aliases_envs,
		Short:   "get solutions template envs",
		Long: "Show list of solution environments." +
			"You can select specific branch specifying branch query (--branch). Default branch is 'master':\n",
		Example: "chkit get template_envs [name]",
		Run: func(cmd *cobra.Command, args []string) {
			logger := coblog.Logger(cmd)
			logger.Debugf("loading solution info")
			if len(args) == 1 {
				var branch, _ = cmd.Flags().GetString("branch")
				var envs, err = ctx.GetClient().GetSolutionsTemplatesEnvs(args[0], branch)
				if err != nil {
					logger.WithError(err).Errorf("unable to get solution list")
					activekit.Attention("Unable to get solution list:\n%v", err)
					ctx.Exit(1)
				}
				if err := export.ExportData(envs, exportConfig); err != nil {
					logrus.WithError(err).Errorf("unable to export data")
					angel.Angel(ctx, err)
				}
			} else {
				cmd.Help()
				ctx.Exit(1)
			}
		},
	}
	command.PersistentFlags().
		String("branch", "", "solution template branch")
	command.PersistentFlags().
		StringVarP((*string)(&exportConfig.Format), "output", "o", "", "output format (yaml/json)")
	command.PersistentFlags().
		StringVarP(&exportConfig.Filename, "file", "f", "", "output file")

	return command
}
