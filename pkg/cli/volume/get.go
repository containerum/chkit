package volume

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/export"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/volume"
	"github.com/ninedraft/boxofstuff/strset"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

var aliases = []string{"volumes", "vol"}

func Get(ctx *context.Context) *cobra.Command {
	var flags struct {
		Output string `flag:"output o" desc:"output format json/yaml"`
		File   string `desc:"output file, STDOUT by default"`
		Names  bool   `desc:"print only names"`
	}
	var command = &cobra.Command{
		Use:     "volume",
		Aliases: aliases,
		Short:   "get volume info",
		Example: "chkit get volume [$VOLUME_NAME...] [-o yaml] [--file $VOLUME_DATA_FILE]",
		Run: func(cmd *cobra.Command, args []string) {
			var logger = ctx.Log.Command("get volumes")
			logger.Debugf("START")
			defer logger.Debugf("END")
			logger.StructFields(flags)
			var renderable model.Renderer
			if len(args) == 1 {
				vol, err := ctx.Client.GetVolume(ctx.Namespace.ID, args[0])
				logger.Debugf("getting volume %q from namespace %q", args[0], ctx.Namespace)
				if err != nil {
					logger.WithError(err).Errorf("unable to get volume %q from namespace %q", args[0], ctx.Namespace)
					fmt.Println(err)
					ctx.Exit(1)
				}
				if flags.Names {
					logger.Debugf("printing name")
					fmt.Println(vol.OwnerAndName())
					return
				}
			} else {
				logger.Debugf("getting volume list from namespace %q", ctx.Namespace)
				list, err := ctx.Client.GetVolumeList(ctx.Namespace.ID)
				if err != nil {
					logger.WithError(err).Errorf("unable to get volume list")
					fmt.Println(err)
					ctx.Exit(1)
				}
				if len(args) > 0 {
					logger.Debugf("filtering volume list by names %v", args)
					var volumeSet = strset.NewSet(args)
					list = list.Filter(func(volume volume.Volume) bool {
						return volumeSet.Have(volume.Name) || volumeSet.Have(volume.OwnerAndName())
					})
				}
				if flags.Names {
					logger.Debugf("printing names")
					fmt.Println(strings.Join(list.OwnersAndNames(), "\n"))
					return
				}
				renderable = list
			}
			logger.Debugf("exporting data")
			if err := export.ExportData(renderable, export.ExportConfig{
				Filename: flags.File,
				Format:   export.ExportFormat(flags.Output),
			}); err != nil {
				logger.WithError(err).Errorf("unable to export data")
				fmt.Println(err)
				ctx.Exit(1)
			}
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
