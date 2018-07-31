package cli

import (
	"fmt"
	"path"
	"strings"
	"unicode/utf8"

	"github.com/blang/semver"
	"github.com/containerum/chkit/help"
	"github.com/containerum/chkit/pkg/cli/doc"
	"github.com/containerum/chkit/pkg/cli/logout"
	"github.com/containerum/chkit/pkg/cli/mode"
	"github.com/containerum/chkit/pkg/cli/postrun"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/cli/set"
	"github.com/containerum/chkit/pkg/cli/setup"
	"github.com/containerum/chkit/pkg/configdir"
	"github.com/containerum/chkit/pkg/context"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

var VERSION = ""

func Root() error {
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
	setup.Config.DebugRequests = true
	if err := setup.SetupLogs(ctx); err != nil {
		return err
	}

	var flags struct {
		Namespace string `flag:"namespace n"`
		Username  string `flag:"username u"`
		Password  string `flag:"password p"`
	}

	prerun.PreRun(ctx, prerun.Config{
		InitClient:         prerun.DoNotInitClient,
		NamespaceSelection: prerun.TemporarySetNamespace,
	})

	root := &cobra.Command{
		Use:   "chkit",
		Short: "Chkit is a terminal client for containerum.io powerful API",
		Long: func() string {
			if ctx.Client.APIaddr == "https://api.containerum.com" {
				var msg = "You are  using Containerum Cloud API. Use 'chkit set api' command to set custom api address"
				var frame = strings.Repeat("!", utf8.RuneCountInString(msg))
				return frame + "\n" + msg + "\n" + frame
			}
			return ""
		}(),
		Version: ctx.Version,
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := prerun.PreRun(ctx, prerun.Config{
				RunLoginOnMissingCreds: true,
				Namespace:              flags.Namespace,
			}); err != nil {
				panic(err)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PostRun: ctx.Defer(func() {
			ctx.Log.Command("root").Debugf("adding postrun")
			postrun.PostRun(ctx)
		}).CobraPostRun,
		TraverseChildren: true,
	}
	ctx.Client.APIaddr = mode.API_ADDR

	if err := gpflag.ParseTo(&flags, root.PersistentFlags()); err != nil {
		panic(err)
	}

	root.AddCommand(RootCommands(ctx)...)
	return root.Execute()
}

func RootCommands(ctx *context.Context) []*cobra.Command {
	var commands = []*cobra.Command{
		setup.Login(ctx),
		Get(ctx),
		Delete(ctx),
		Create(ctx),
		Replace(ctx),
		set.Set(ctx),
		Logs(ctx),
		Run(ctx),
		Rename(ctx),
		logout.Logout(ctx),
		Update(ctx),
		{
			Use:   "version",
			Short: "Print version",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(ctx.Version)
			},
		},
		doc.Doc(ctx),
	}
	help.AutoForCommands(commands)
	return commands
}

func RootCommandsWithEmptyContext() []*cobra.Command {
	var ctx = &context.Context{}
	var commands = []*cobra.Command{
		setup.Login(ctx),
		Get(ctx),
		Delete(ctx),
		Create(ctx),
		Replace(ctx),
		set.Set(ctx),
		Logs(ctx),
		Run(ctx),
		Rename(ctx),
		logout.Logout(ctx),
		Update(ctx),
		{
			Use:   "version",
			Short: "Print version",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(ctx.Version)
			},
		},
		doc.Doc(ctx),
	}
	help.AutoForCommands(commands)
	return commands
}
