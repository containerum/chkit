package clipod

import (
	"fmt"

	"os"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/export"
	"github.com/containerum/chkit/pkg/model/pod"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/strset"
	"github.com/spf13/cobra"
)

var aliases = []string{"po", "pods"}

var getPodConfig = struct {
	export.ExportConfig
}{
	ExportConfig: export.ExportConfig{
		Filename: "-",
	},
}

func Get(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:     "pod",
		Aliases: aliases,
		Short:   "show pod info",
		Long:    "Show pod info.",
		Example: "chkit get pod pod_label [-o yaml/json] [-f output_file]",
		Run: func(cmd *cobra.Command, args []string) {
			switch len(args) {
			case 0:
				polist, err := ctx.Client.GetPodList(ctx.Namespace.ID)
				if err != nil {
					fmt.Printf("Unable to get pod list from namespace %q :(\n", ctx.Namespace)
					os.Exit(1)
				}
				if err := export.ExportData(polist, getPodConfig.ExportConfig); err != nil {
					activekit.Attention(err.Error())
					os.Exit(1)
				}
				return
			case 1:
				po, err := ctx.Client.GetPod(ctx.Namespace.ID, args[0])
				if err != nil {
					fmt.Printf("Unable to get pod from namespace %q :(\n", ctx.Namespace)
					os.Exit(1)
				}
				if err := export.ExportData(po, getPodConfig.ExportConfig); err != nil {
					activekit.Attention(err.Error())
					os.Exit(1)
				}
			default:
				polist, err := ctx.Client.GetPodList(ctx.Namespace.ID)
				if err != nil {
					fmt.Printf("Unable to get pod list from namespace %q :(\n", ctx.Namespace)
					os.Exit(1)
				}
				var filteredList pod.PodList = make([]pod.Pod, 0, len(polist))
				podNames := strset.NewSet(args)
				for _, p := range polist {
					if podNames.Have(p.Name) {
						filteredList = append(filteredList, p)
					}
				}
				if err := export.ExportData(filteredList, getPodConfig.ExportConfig); err != nil {
					activekit.Attention(err.Error())
					os.Exit(1)
				}
				return
			}
		},
	}
	command.PersistentFlags().
		StringVarP((*string)(&getPodConfig.Format), "output", "o", "", "output format (json/yaml)")
	command.PersistentFlags().
		StringVarP(&getPodConfig.Filename, "file", "f", "-", "output file")
	return command
}
