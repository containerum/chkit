package clipod

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
	podControl "github.com/containerum/chkit/pkg/controls/pod"
	"github.com/containerum/chkit/pkg/export"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

var aliases = []string{"po", "pods"}

func Get(ctx *context.Context) *cobra.Command {
	var flags podControl.GetFlags
	command := &cobra.Command{
		Use:     "pod",
		Aliases: aliases,
		Short:   "show pod info",
		Long:    "Show pod info.",
		Example: "chkit get pod pod_label [-o yaml/json] [-f output_file]",
		Run: func(cmd *cobra.Command, args []string) {
			var logger = ctx.Log.Command("get pod")
			logger.Debugf("START")
			defer logger.Debugf("END")
			logger.StructFields(flags)
			switch len(args) {
			case 1:
				logger.Debugf("getting pod %q from namespace %q", args[0], ctx.GetNamespace())
				po, err := ctx.GetClient().GetPod(ctx.GetNamespace().ID, args[0])
				if err != nil {
					logger.WithError(err).Errorf("unable to get pod %q from namespace %q", args[0], ctx.GetNamespace())
					fmt.Printf("Unable to get pod from namespace %q :(\n", ctx.GetNamespace())
					ctx.Exit(1)
				}
				logger.Debugf("exporting data")
				if err := export.ExportData(po, flags.ExportConfig()); err != nil {
					activekit.Attention(err.Error())
					ctx.Exit(1)
				}
			default:
				logger.Debugf("getting pod list from namespace %q", ctx.GetNamespace())
				polist, err := ctx.GetClient().GetPodList(ctx.GetNamespace().ID)
				if err != nil {
					logger.WithError(err).Errorf("unable to get deployment list from namespace %q", ctx.GetNamespace())
					fmt.Printf("Unable to get pod list from namespace %q :(\n", ctx.GetNamespace())
					ctx.Exit(1)
				}
				if flags.IsStatusesDefined() {
					logger.Debugf("filtering pod list by statuses %v", flags.Statuses())
					polist = polist.Filter(flags.StatusFilter())
				}
				if len(args) > 0 {
					logger.Debugf("selecting pod with names %v", args)
					polist = polist.FilterByNames(args...)
				}
				if err := export.ExportData(polist, flags.ExportConfig()); err != nil {
					activekit.Attention(err.Error())
					ctx.Exit(1)
				}
				return
			}
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
