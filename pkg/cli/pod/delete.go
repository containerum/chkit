package clipod

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/export"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/model/pod"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Delete(ctx *context.Context) *cobra.Command {
	var deletePodConfig = struct {
		Force bool
	}{}
	exportConfig := export.ExportConfig{}
	command := &cobra.Command{
		Use:     "pod",
		Aliases: aliases,
		Short:   "delete pod in specific namespace",
		Long:    "Delete pods.",
		Example: "chkit delete pod pod_name [-n namespace]",
		Run: func(cmd *cobra.Command, args []string) {
			var selectedPod string
			var selectedDeploy string
			if len(args) == 0 && !deletePodConfig.Force {
				deployList, err := ctx.Client.GetDeploymentList(ctx.GetNamespace().ID)
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				var deployMenu activekit.MenuItems
				for _, depl := range deployList {
					deployMenu = deployMenu.Append(&activekit.MenuItem{
						Label: depl.Name,
						Action: func(depl deployment.Deployment) func() error {
							return func() error {
								selectedDeploy = depl.Name
								return nil
							}
						}(depl.Copy()),
					})
				}
				(&activekit.Menu{
					Title: "Select deployment",
					Items: deployMenu,
				}).Run()

				podList, err := ctx.Client.GetDeploymentPodList(ctx.GetNamespace().ID, selectedDeploy)
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				if err := export.ExportData(podList, exportConfig); err != nil {
					logrus.WithError(err).Errorf("unable to export data")
					angel.Angel(ctx, err)
				}
				var menu activekit.MenuItems
				for _, pd := range podList {
					menu = menu.Append(&activekit.MenuItem{
						Label: pd.Name,
						Action: func(pd pod.Pod) func() error {
							return func() error {
								selectedPod = pd.Name
								return nil
							}
						}(pd.Copy()),
					})
				}
				(&activekit.Menu{
					Title: "Select pod",
					Items: menu,
				}).Run()
			} else {
				if len(args) == 0 {
					cmd.Help()
					ctx.Exit(1)
				}
				selectedPod = args[0]
			}
			if deletePodConfig.Force || activekit.YesNo("Are you sure you want to delete pod %q in namespace %q?", selectedPod, ctx.GetNamespace()) {
				if err := ctx.Client.DeletePod(ctx.GetNamespace().ID, selectedPod); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Println("OK")
			}
		},
	}
	command.PersistentFlags().
		BoolVarP(&deletePodConfig.Force, "force", "f", false, "delete pod without confirmation")
	return command
}
