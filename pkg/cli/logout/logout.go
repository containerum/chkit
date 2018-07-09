package logout

import (
	"os"
	"path"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/ninedraft/boxofstuff/str"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Logout(ctx *context.Context) *cobra.Command {
	var flags struct {
		Purge bool `desc:"clean all config files, logs and reports"`
		Force bool
	}
	var command = &cobra.Command{
		Use:   "logout",
		Short: "Logout from chkit, delete garbage files",
		Run: func(cmd *cobra.Command, args []string) {
			var filesToRemove = str.Vector{"tokens", "logs", "reports"}
			if flags.Purge {
				filesToRemove = append(filesToRemove, "config.toml")
			}
			if flags.Force || activekit.YesNo("The following files will be removed:\n%s\n"+
				"Are sure you want to logout from chkit?", filesToRemove.Join("\n")) {
				switch flags.Purge {
				case true:
					if err := os.RemoveAll(ctx.ConfigPath); err != nil {
						ferr.Println(err)
						ctx.Exit(1)
					}
				case false:
					defer func() { ctx.Changed = true }()
					ctx.Client.UserInfo = model.UserInfo{}
					if err := os.Remove(path.Join(ctx.ConfigDir, "tokens")); err != nil {
						ferr.Println(err)
					}
					if err := os.RemoveAll(path.Join(ctx.ConfigDir, "support")); err != nil {
						ferr.Println(err)
					}
				}
			}
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
