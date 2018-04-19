package replicas

import (
	"github.com/containerum/chkit/pkg/context"
	"github.com/spf13/cobra"
	"strconv"
	"github.com/containerum/chkit/pkg/util/activekit"
	"fmt"
	"os"
	"github.com/containerum/chkit/pkg/model/deployment"
)

func Set(ctx *context.Context) *cobra.Command {
	var deplName string
	var replicas uint64
	command := &cobra.Command{
		Use: "replicas",
		Short:"set deployment replicas",
		Long: "Sets deployment replicas",
		Example: "chkit set replicas [-n namespace_label] [-d depl_label] [N_replicas]",
		Aliases: []string{"re", "rep", "repl", "replica"},
		Run: func(cmd *cobra.Command, args []string) {
			if !cmd.Flag("deployment").Changed {
				deplList, err := ctx.Client.GetDeploymentList(ctx.Namespace)
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
					Items:menu,
				}).Run()
			}
			if !cmd.Flag("replicas").Changed {
				if len(args) > 1 {
					cmd.Help()
					return
				}
				n, err := strconv.ParseUint(args[0], 10, 64)
				if err != nil {
					activekit.Attention(fmt.Sprintf("replicas parameter must be number 1..15:\n%v\n", err))
					os.Exit(1)
				}
				replicas = n
			}
			if replicas < 1 || replicas > 15 {
				activekit.Attention(fmt.Sprintf("replicas parameter must be number 1..15, but it %d\n", replicas))
				os.Exit(1)
			}
		if err := ctx.Client.SetReplicas(ctx.Namespace, deplName, replicas); err != nil {
			activekit.Attention(err.Error())
			os.Exit(1)
		}
		fmt.Println("OK")
		},
	}
	command.PersistentFlags().
		StringVarP(&deplName, "deployment", "d", "", "deployment name")
	command.PersistentFlags() .
		Uint64VarP(&replicas, "replicas", "r", 1, "replicas, 1..15")
	return command
}
