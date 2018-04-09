package cmd

import (
	"github.com/containerum/chkit/pkg/update"
	"gopkg.in/urfave/cli.v2"
)

var commandUpdate = &cli.Command{
	Name:      "update",
	Usage:     "update chkit client",
	UsageText: "chkit update [from github|dir]",
	Action: func(ctx *cli.Context) error {
		return updateFromGithub(ctx)
	},
	Subcommands: []*cli.Command{
		&cli.Command{
			Name: "from",
			Subcommands: []*cli.Command{
				&cli.Command{
					Name:   "github",
					Action: updateFromGithub,
				},
				&cli.Command{
					Name:   "dir",
					Action: updateFromDir,
				},
			},
		},
	},
}

func updateFromGithub(ctx *cli.Context) error {
	return update.Update(ctx,
		update.NewGithubLatestCheckerDownloader(ctx,
			"containerum",
			"chkit"), false)
}

func updateFromDir(ctx *cli.Context) error {
	return update.Update(ctx,
		update.NewFileSystemLatestCheckerDownloader(ctx,
			ctx.Args().First()), false)
}
