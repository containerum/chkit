package clideployment

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Delete(ctx *context.Context) *cobra.Command {
	var deleteDeplConfig = struct {
		Force bool
	}{}
	command := &cobra.Command{
		Use:     "deployment",
		Aliases: aliases,
		Short:   "call to delete deployment in specific namespace",
		Long: `Deletes deployment in specific namespace.
Use --force flag to suppress confirmation`,
		Run: func(cmd *cobra.Command, args []string) {
			switch len(args) {
			case 0:
				list, err := ctx.Client.GetDeploymentList(ctx.Namespace.ID)
				if err != nil {
					logrus.WithError(err).Errorf("unable to get deployment list")
					activekit.Attention(err.Error())
				}
				var menu []*activekit.MenuItem
				for _, depl := range list {
					menu = append(menu, &activekit.MenuItem{
						Label: depl.Name,
						Action: func(depl deployment.Deployment) func() error {
							return func() error {
								if activekit.YesNo(fmt.Sprintf("Are you sure you want to delete deployment %q?", depl.Name)) {
									if err := ctx.Client.DeleteDeployment(ctx.Namespace.ID, depl.Name); err != nil {
										logrus.WithError(err).Debugf("unable to delete deployment %q in namespace %q", depl.Name, ctx.Namespace)
										activekit.Attention(err.Error())
										os.Exit(1)
									}
									fmt.Printf("OK\n")
									logrus.Debugf("deployment %q in namespace %q deleted", depl.Name, ctx.Namespace)
								}
								return nil
							}
						}(depl),
					})
				}
				(&activekit.Menu{
					Items: append(menu, []*activekit.MenuItem{
						{
							Label: "Exit",
						},
					}...),
				}).Run()
			case 1:
				deplName := args[0]
				logrus.
					WithField("command", "delete deployment").
					Debugf("start deleting deployment %q", deplName)
				if deleteDeplConfig.Force || activekit.YesNo(fmt.Sprintf("Are you sure you want to delete deployment %q?", deplName)) {
					if err := ctx.Client.DeleteDeployment(ctx.Namespace.ID, deplName); err != nil {
						logrus.WithError(err).Debugf("unable to delete deployment %q in namespace %q", deplName, ctx.Namespace)
						activekit.Attention(err.Error())
						os.Exit(1)
					}
					fmt.Printf("OK\n")
					logrus.Debugf("deployment %q in namespace %q deleted", deplName, ctx.Namespace)
				}
			default:
				cmd.Help()
				return
			}

		},
	}
	command.PersistentFlags().
		BoolVarP(&deleteDeplConfig.Force, "force", "f", false, "delete without confirmation")
	return command
}
