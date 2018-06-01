package clipod

import (
	"strings"

	"fmt"

	"os"

	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/pod"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/strset"
	"github.com/spf13/cobra"
)

var aliases = []string{"po", "pods"}

var getPodConfig = struct {
	configuration.ExportConfig
}{
	ExportConfig: configuration.ExportConfig{
		Filename: "-",
	},
}

func Get(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:     "pod",
		Aliases: aliases,
		Short:   "shows pod info",
		Long:    "shows pod info. Aliases: " + strings.Join(aliases, ", "),
		Example: "chkit get pod pod_label [-o yaml/json] [-f output_file]",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := prerun.PreRun(ctx); err != nil {
				angel.Angel(ctx, err)
				os.Exit(1)
			}
			if cmd.Flags().Changed("namespace") {
				ctx.Namespace.ID, _ = cmd.Flags().GetString("namespace")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			switch len(args) {
			case 0:
				polist, err := ctx.Client.GetPodList(ctx.Namespace.ID)
				if err != nil {
					fmt.Printf("Unable to get pod list from namespace %q :(\n", ctx.Namespace)
					os.Exit(1)
				}
				if err := configuration.ExportData(polist, getPodConfig.ExportConfig); err != nil {
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
				if err := configuration.ExportData(po, getPodConfig.ExportConfig); err != nil {
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
				if err := configuration.ExportData(filteredList, getPodConfig.ExportConfig); err != nil {
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
