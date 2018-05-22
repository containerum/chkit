package cli

import (
	"os"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/update"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/spf13/cobra"
)

func Update(ctx *context.Context) *cobra.Command {
	var debug bool
	command := &cobra.Command{
		Use:     "update",
		Short:   "update chkit client",
		Example: "chkit update [from github|dir <path>] [--debug]",
		Run: func(cmd *cobra.Command, args []string) {
			if err := updateFromGithub(debug); err != nil {
				activekit.Attention(err.Error())
				os.Exit(1)
			}
		},
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
		Run: func(cmd *cobra.Command, args []string) {
			if err := updateFromGithub(*debug); err != nil {
				activekit.Attention(err.Error())
				os.Exit(1)
			}
		},
	}
	return command
}

func updateFromDirCommand(ctx *context.Context, debug *bool) *cobra.Command {
	command := &cobra.Command{
		Use:     "dir",
		Short:   "update from local directory",
		Example: "chkit update from dir <path> [--debug]",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) <= 0 {
				cmd.Help()
				os.Exit(1)
			}
			if err := updateFromDir(args[0]); err != nil {
				activekit.Attention(err.Error())
				os.Exit(1)
			}
		},
	}
	return command
}

func updateFromGithub(debug bool) error {
	return update.Update(update.NewGithubLatestCheckerDownloader("containerum", "chkit", debug), false)
}

func updateFromDir(path string) error {
	return update.Update(update.NewFileSystemLatestCheckerDownloader(path), false)
}
