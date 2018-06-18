package replicas

import (
	"fmt"
	"os"
	"strconv"

	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/spf13/cobra"
)

func Set(ctx *context.Context) *cobra.Command {
	var deplName string
	var replicas uint64
	command := &cobra.Command{
		Use:     "replicas",
		Short:   "Set deployment replicas",
		Long:    "Set deployment replicas.",
		Example: "chkit set replicas [-n namespace_label] [-d depl_label] [N_replicas]",
		Aliases: []string{"re", "rep", "repl", "replica"},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := prerun.PreRun(ctx); err != nil {
				activekit.Attention(err.Error())
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if !cmd.Flag("deployment").Changed {
				deplList, err := ctx.Client.GetDeploymentList(ctx.Namespace.ID)
				if err != nil {
					activekit.Attention(fmt.Sprintf("Unable to get deployment list:\n%v", err))
					os.Exit(1)
				}
				var menu []*activekit.MenuItem
				for _, depl := range deplList {
					menu = append(menu, &activekit.MenuItem{
						Label: depl.Name,
						Action: func(depl deployment.Deployment) func() error {
							return func() error {
								deplName = depl.Name
								return nil
							}
						}(depl),
					})
				}
				(&activekit.Menu{
					Title: "Select deployment",
					Items: menu,
				}).Run()
			}
			if !cmd.Flag("replicas").Changed {
				for {
					n, err := strconv.ParseUint(activekit.Promt("Type replicas number: "), 10, 64)
					if err != nil {
						activekit.Attention(fmt.Sprintf("replicas parameter must be number 1..15:\n%v\n", err))
						continue
					}
					replicas = n
					break
				}
			}
			if replicas < 1 || replicas > 15 {
				activekit.Attention(fmt.Sprintf("replicas parameter must be number 1..15, but it %d\n", replicas))
				os.Exit(1)
			}
			if err := ctx.Client.SetReplicas(ctx.Namespace.ID, deplName, replicas); err != nil {
				activekit.Attention(err.Error())
				os.Exit(1)
			}
			fmt.Println("OK")
		},
	}
	command.PersistentFlags().
		StringVarP(&deplName, "deployment", "d", "", "deployment name")
	command.PersistentFlags().
		Uint64VarP(&replicas, "replicas", "r", 1, "replicas, 1..15")
	return command
}
