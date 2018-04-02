package cmd

import (
	"io"
	"strings"

	"os"

	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/client"
	"gopkg.in/urfave/cli.v2"
)

var logsCommandAliases = []string{"log"}
var commandLogs = &cli.Command{
	Name:        "logs",
	Aliases:     logsCommandAliases,
	Description: `View pod logs`,
	Usage:       `view pod logs. Aliases: ` + strings.Join(logsCommandAliases, ", "),
	UsageText:   `logs pod_label [container] [--follow] [--prev] [--tail n]`,
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
		&cli.BoolFlag{
			Name:    "prev",
			Aliases: []string{"p"},
			Usage:   `show logs from previous instance (useful for crashes debugging)`,
		},
		&cli.IntFlag{
			Name:    "tail",
			Aliases: []string{"t"},
			Value:   100,
			Usage:   `print last <value> log lines`,
		},
		util.NamespaceFlag,
	},
}
