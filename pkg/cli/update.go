package cli

import (
	"github.com/blang/semver"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/update"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/spf13/cobra"
)

func Update(ctx *context.Context) *cobra.Command {
	var debug bool
	command := &cobra.Command{
		Use:     "update",
		Short:   "Update chkit client",
		Long:    `Use "chkit update [command] --help" for more information about the command.`,
		Example: "chkit update [from github|dir <path>] [--debug]",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := prerun.PreRun(ctx); err != nil {
				angel.Angel(ctx, err)
				ctx.Exit(1)
			}
			if err := prerun.GetNamespaceByUserfriendlyID(ctx, cmd.Flags()); err != nil {
				ferr.Println(err)
				ctx.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := updateFromGithub(ctx, debug); err != nil {
				activekit.Attention(err.Error())
			}
		},
		PersistentPostRun: ctx.CobraPostRun,
	}
	command.PersistentFlags().
		BoolVarP(&debug, "debug", "", false, "print debug information")
	command.AddCommand(fromCommand(ctx, &debug))
	return command
}

func fromCommand(ctx *context.Context, debug *bool) *cobra.Command {
	command := &cobra.Command{
		Use: "from",
	}
	command.AddCommand(updateFromGithubCommand(ctx, debug))
	command.AddCommand(updateFromDirCommand(ctx, debug))
	return command
}

func updateFromGithubCommand(ctx *context.Context, debug *bool) *cobra.Command {
	command := &cobra.Command{
		Use:   "github",
		Short: "update from github releases",
		Long:  "Update from github releases.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := updateFromGithub(ctx, *debug); err != nil {
				activekit.Attention(err.Error())
				ctx.Exit(1)
			}
		},
	}
	return command
}

func updateFromDirCommand(ctx *context.Context, debug *bool) *cobra.Command {
	command := &cobra.Command{
		Use:     "dir",
		Short:   "update from local directory",
		Long:    "Update from local directory.",
		Example: "chkit update from dir <path> [--debug]",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) <= 0 {
				cmd.Help()
				ctx.Exit(1)
			}
			if err := updateFromDir(ctx, args[0]); err != nil {
				activekit.Attention(err.Error())
				ctx.Exit(1)
			}
		},
	}
	return command
}

func updateFromGithub(ctx *context.Context, debug bool) error {
	ver, err := semver.ParseTolerant(ctx.Version)
	if err != nil {
		return err
	}
	return update.Update(
		ver,
		update.NewGithubLatestCheckerDownloader("containerum", "chkit", debug),
		false,
	)
}

func updateFromDir(ctx *context.Context, path string) error {
	ver, err := semver.ParseTolerant(ctx.Version)
	if err != nil {
		return err
	}
	return update.Update(
		ver,
		update.NewFileSystemLatestCheckerDownloader(path),
		false,
	)
}
