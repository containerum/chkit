package cli

import (
	"fmt"
	"path"

	"github.com/blang/semver"
	"github.com/containerum/chkit/pkg/cli/doc"
	"github.com/containerum/chkit/pkg/cli/mode"
	"github.com/containerum/chkit/pkg/cli/postrun"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/cli/set"
	"github.com/containerum/chkit/pkg/cli/setup"
	"github.com/containerum/chkit/pkg/configdir"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

var VERSION = ""

func Root() error {
	var flags struct {
		Username  string
		Password  string
		Namespace string
	}
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
	setup.SetupLogs(ctx)

	root := &cobra.Command{
		Use:     "chkit",
		Short:   "Chkit is a terminal client for containerum.io powerful API",
		Version: ctx.Version,
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := prerun.PreRun(ctx, prerun.Config{
				RunLoginOnMissingCreds: true,
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
		}).CobraPostrun,
		TraverseChildren: true,
	}
	ctx.GetClient().APIaddr = mode.API_ADDR
	if err := gpflag.ParseTo(&flags, root.PersistentFlags()); err != nil {
		ferr.Println(err)
		ctx.Exit(1)
	}

	root.AddCommand(
		setup.Login(ctx),
		Get(ctx),
		Delete(ctx),
		Create(ctx),
		Replace(ctx),
		set.Set(ctx),
		Logs(ctx),
		Run(ctx),
		Rename(ctx),
		Update(ctx),
		&cobra.Command{
			Use:   "version",
			Short: "Print version",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(ctx.Version)
			},
		},
		doc.Doc(ctx),
	)
	return root.Execute()
}

func RootCommands() []*cobra.Command {
	var ctx = &context.Context{}
	return []*cobra.Command{
		setup.Login(ctx),
		Get(ctx),
		Delete(ctx),
		Create(ctx),
		Replace(ctx),
		set.Set(ctx),
		Logs(ctx),
		Run(ctx),
		Rename(ctx),
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
}
