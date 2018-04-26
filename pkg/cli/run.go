package cli

import (
	"github.com/containerum/chkit/pkg/cli/mode"

	"path"

	"fmt"

	"os"

	"github.com/blang/semver"
	"github.com/containerum/chkit/pkg/cli/clisetup"
	"github.com/containerum/chkit/pkg/cli/login"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/cli/set"
	"github.com/containerum/chkit/pkg/configdir"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var VERSION = ""

func Run() error {
	ctx := &context.Context{
		Version: func() string {
			// try to normalise version string
			v, err := semver.ParseTolerant(VERSION)
			if err != nil {
				return VERSION
			}
			return v.String()
		}(),
		ConfigDir:  configdir.ConfigDir(),
		ConfigPath: path.Join(configdir.ConfigDir(), "config.toml"),
	}

	root := &cobra.Command{
		Use:     "chkit",
		Short:   "Chkit is a terminal client for containerum.io powerful API",
		Version: ctx.Version,
		PreRun: func(cmd *cobra.Command, args []string) {
			clisetup.Config.DebugRequests = true
			clisetup.SetupLogs(ctx)
			if cmd.Flag("username").Changed && cmd.Flag("password").Changed {
				if err := login.Setup(ctx); err != nil {
					angel.Angel(ctx, err)
					os.Exit(1)
				}
				return
			} else if cmd.Flag("username").Changed || cmd.Flag("password").Changed {
				cmd.Help()
				os.Exit(1)
			}
			if err := prerun.PreRun(ctx); err != nil {
				angel.Angel(ctx, err)
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Hello, %q!\nUsing %q as default namespace\n",
				ctx.Client.Username,
				ctx.Namespace)
			if err := mainActivity(); err != nil {
				logrus.Fatalf("error in main activity: %v", err)
			}
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			if !ctx.Changed {
				return
			}
			if err := configuration.SaveConfig(ctx); err != nil {
				fmt.Printf("Unable to save config file: %v\n", err)
			}
		},
		TraverseChildren: true,
	}
	ctx.Client.APIaddr = mode.API_ADDR

	root.PersistentFlags().
		StringVarP(&ctx.Client.Username, "username", "u", "", "account username")
	root.PersistentFlags().
		StringVarP(&ctx.Client.Password, "password", "p", "", "account password")
	root.PersistentFlags().
		StringVarP(&ctx.Namespace, "namespace", "n", ctx.Namespace, "")
	root.PersistentFlags().
		BoolVarP(&ctx.Quiet, "quiet", "q", ctx.Quiet, "quiet mode")

	root.AddCommand(
		login.Login(ctx),
		Get(ctx),
		Delete(ctx),
		Create(ctx),
		Replace(ctx),
		set.Set(ctx),
		//Logs(ctx),
	)
	return root.Execute()
}
