package cmd

import (
	"io"

	"os"

	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/client"
	"gopkg.in/urfave/cli.v2"
)

var commandLogs = &cli.Command{
	Name:        "logs",
	Description: `View pod logs`,
	Usage:       `view pod logs`,
	UsageText:   `logs [command options] <pod name> [container name]`,
	Before: func(ctx *cli.Context) error {
		if ctx.Bool("help") {
			return cli.ShowSubcommandHelp(ctx)
		}
		return setupAll(ctx)
	},
	Action: func(ctx *cli.Context) error {
		client := util.GetClient(ctx)
		defer util.StoreClient(ctx, client)
		var podName string
		var containerName string
		switch ctx.NArg() {
		case 2:
			containerName = ctx.Args().Tail()[0]
			fallthrough
		case 1:
			podName = ctx.Args().First()
		default:
			cli.ShowSubcommandHelp(ctx)
			return nil
		}

		params := chClient.GetPodLogsParams{
			Namespace: util.GetNamespace(ctx),
			Pod:       podName,
			Container: containerName,
			Follow:    ctx.Bool("follow"),
			Previous:  ctx.Bool("previous"),
			Tail:      ctx.Int("tail"),
		}
		rc, err := client.GetPodLogs(params)
		if err != nil {
			return err
		}
		defer rc.Close()

		io.Copy(os.Stdout, rc)

		return nil
	},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "follow",
			Aliases: []string{"f"},
			Usage:   `follow pod logs`,
		},
		&cli.StringFlag{
			Name:    "prev",
			Aliases: []string{"p"},
			Usage:   `show logs from previous instance (useful for crashes debugging)`,
		},
		&cli.IntFlag{
			Name:    "tail",
			Aliases: []string{"t"},
			Usage:   `print last <value> log lines`,
		},
		util.NamespaceFlag,
	},
}
