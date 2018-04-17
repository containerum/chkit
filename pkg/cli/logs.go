package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	ErrUnableToReadLogs chkitErrors.Err = "unable to read logs"
)

var logsCommandAliases = []string{"log"}

func Logs(ctx *context.Context) *cobra.Command {
	var logsConfig = struct {
		Quiet  bool
		Follow bool
		Prev   bool
		Tail   uint
	}{}
	command := &cobra.Command{
		Use:     "logs",
		Aliases: logsCommandAliases,
		Short:   "View pod logs",
		Long:    `view pod logs. Aliases: ` + strings.Join(logsCommandAliases, ", "),
		Example: `logs pod_label [container] [--follow] [--prev] [--tail n] [--quiet]`,
		Run: func(cmd *cobra.Command, args []string) {
			var podName string
			var containerName string
			client := context.GlobalContext.Client
			switch len(args) {
			case 2:
				containerName = args[1]
				fallthrough
			case 1:
				podName = args[0]
			default:
				cmd.Help()
				return
			}

			params := chClient.GetPodLogsParams{
				Namespace: context.GlobalContext.Namespace,
				Pod:       podName,
				Container: containerName,
				Follow:    logsConfig.Follow,
				Previous:  logsConfig.Prev,
				Tail:      int(logsConfig.Tail),
			}
			rc, err := client.GetPodLogs(params)
			if err != nil {
				logrus.WithError(err).Errorf("error while getting logs")
				activekit.Attention(err.Error())
			}
			defer rc.Close()

			scanner := bufio.NewScanner(rc)
			var nLines uint64
			for scanner.Scan() {
				if err := scanner.Err(); err != nil {
					err = ErrUnableToReadLogs.Wrap(err)
					logrus.WithError(err).Errorf("unable to scan logs byte stream")
					activekit.Attention(err.Error())
				}
				fmt.Println(scanner.Text())
				nLines++
			}
			fmt.Printf("%d lines of logs read\n", nLines)
			activekit.Attention(err.Error())
		},
	}
	command.PersistentFlags().
		BoolVarP(&logsConfig.Quiet, "quiet", "q", false, "print only logs and errors")
	command.PersistentFlags().
		BoolVarP(&logsConfig.Follow, "follow", "f", false, `follow pod logs`)
	command.PersistentFlags().
		UintVarP(&logsConfig.Tail, "tail", "t", 100, `print last <value> log lines`)
	return command
}
