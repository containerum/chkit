package clitemplate

import (
	"os"

	"fmt"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/solution"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/spf13/cobra"
)

var aliases_envs = []string{"template_env", "tmpl_env", "envs", "environments", "templates_environments", "tmpls_env", "tmpenv", "tmp_env", "tmps_env", "tmpsenv"}

func GetEnvs(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:     "template_envs",
		Aliases: aliases_envs,
		Short:   "get solutions template envs",
		Long:    "Show list of solution enviroments.",
		Example: "chkit get template_envs [name]",
		Run: func(cmd *cobra.Command, args []string) {
			logger := coblog.Logger(cmd)
			logger.Debugf("loading solution info")
			if len(args) == 1 {
				var branch, _ = cmd.Flags().GetString("branch")
				var envs, err = ctx.Client.GetSolutionsTemplatesEnvs(args[0], branch)
				if err != nil {
					logger.WithError(err).Errorf("unable to get solution list")
					activekit.Attention("Unable to get solution list:\n%v", err)
					os.Exit(1)
				}
				fmt.Println(solution.SolutionEnvFromKube(envs).RenderTable())
			} else {
				cmd.Help()
				os.Exit(1)
			}
		},
	}
	command.PersistentFlags().
		String("branch", "", "branch")
	return command
}
