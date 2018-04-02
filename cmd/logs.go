package cmd

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/sirupsen/logrus"

	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/client"
	"gopkg.in/urfave/cli.v2"
)

const (
	ErrUnableToReadLogs chkitErrors.Err = "unable to read logs"
)

var logsCommandAliases = []string{"log"}
var commandLogs = &cli.Command{
	Name:        "logs",
	Aliases:     logsCommandAliases,
	Description: `View pod logs`,
	Usage:       `view pod logs. Aliases: ` + strings.Join(logsCommandAliases, ", "),
	UsageText:   `logs pod_label [container] [--follow] [--prev] [--tail n] [--quiet]`,
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
			logrus.WithError(err).Errorf("error while getting logs")
			return err
		}
		defer rc.Close()

		scanner := bufio.NewScanner(rc)
		var nLines uint64
		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				err = ErrUnableToReadLogs.Wrap(err)
				logrus.WithError(err).Errorf("unable to scan logs byte stream")
				return err
			}
			fmt.Println(scanner.Text())
			nLines++
		}
		fmt.Printf("%d lines of logs read\n", nLines)
		return nil
	},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "quiet",
			Aliases: []string{"q"},
			Usage:   "print only logs and errors",
			Value:   false,
		},
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
