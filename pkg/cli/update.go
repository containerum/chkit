package cli

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
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:   "debug-requests",
			Hidden: true,
		},
	},
	Subcommands: []*cli.Command{
		{
			Name: "from",
			Subcommands: []*cli.Command{
				{
					Name:   "github",
					Action: updateFromGithub,
				},
				{
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
