package volume

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Delete(ctx *context.Context) *cobra.Command {
	var flags struct {
		Force bool `desc:"suppress confirmation" flag:"force f"`
	}
	var command = &cobra.Command{
		Use:     "delete",
		Short:   "delete volume",
		Example: "chkit delete volume [--force]",
		Aliases: aliases,
		Run: func(cmd *cobra.Command, args []string) {
			var logger = ctx.Log.Command("delete volume")
			logger.Debugf("START")
			defer logger.Debugf("END")
			logger.StructFields(flags)
			var volumeID string
			logger.Debugf("getting volume list")
			var volumeList, err = ctx.Client.GetVolumeList(ctx.Namespace.ID)
			if err != nil {
				logger.WithError(err).Errorf("unable to get volume list")
				fmt.Println(err)
				os.Exit(1)
			}
			switch len(args) {
			case 0:
				if flags.Force {
					fmt.Println("if flag --force is active then volume name must be provided as first arg")
					os.Exit(1)
				}
				logger.Debugf("selecting volume in interactive mode")
				(&activekit.Menu{
					Title: "Select volume",
					Items: activekit.StringSelector(volumeList.Names(), func(s string) error {
						volumeID = s
						logger.Debugf("volume %q selected", s)
						return nil
					}),
				}).Run()
			case 1:
				logger.Debugf("searching volume %q", args[0])
				var vol, ok = volumeList.GetByUserFriendlyID(args[0])
				if !ok {
					logger.Debugf("volume %q not found", args[0])
					fmt.Printf("volume %q not found!\n", args[0])
					os.Exit(1)
				}
				logger.Debugf("found volume %q", args[0])
				volumeID = vol.Name
			default:
				cmd.Help()
				os.Exit(1)
			}
			if flags.Force || activekit.YesNo("Do you really want to delete volume %q?", volumeID) {
				logger.Debugf("deleting volume %q in namespace %q", volumeID, ctx.Namespace)
				if err := ctx.Client.DeleteVolume(ctx.Namespace.ID, volumeID); err != nil {
					logger.WithError(err).Errorf("unable to delete volume %q in namespace %q", volumeID, ctx.Namespace)
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println("OK")
			}
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
