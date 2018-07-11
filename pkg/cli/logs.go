package cli

import (
	"bufio"
	"fmt"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	ErrUnableToReadLogs chkitErrors.Err = "unable to read logs"
)

var logsCommandAliases = []string{"log"}

func Logs(ctx *context.Context) *cobra.Command {
	var logsConfig = struct {
		Follow bool
		Prev   bool
		Tail   uint
	}{}
	command := &cobra.Command{
		Use:     "logs",
		Aliases: logsCommandAliases,
		Short:   "View pod logs",
		Example: `logs pod_label [container] [--follow] [--prev] [--tail n] [--quiet]`,
		PreRun: func(cmd *cobra.Command, args []string) {
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
			var podName string
			var containerName string
			client := ctx.GetClient()
			switch len(args) {
			case 2:
				containerName = args[1]
				fallthrough
			case 1:
				podName = args[0]
			default:
				var pods, err = client.GetPodList(ctx.GetNamespace().ID)
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				(&activekit.Menu{
					Title: "Select pod",
					Items: activekit.ItemsFromIter(uint(pods.Len()), func(index uint) *activekit.MenuItem {
						var po = pods[index]
						return &activekit.MenuItem{
							Label: po.Name,
							Action: func() error {
								podName = po.Name
								return nil
							},
						}
					}),
				}).Run()
			}

			params := chClient.GetPodLogsParams{
				Namespace: ctx.GetNamespace().ID,
				Pod:       podName,
				Container: containerName,
				Follow:    logsConfig.Follow,
				Previous:  logsConfig.Prev,
				Tail:      int(logsConfig.Tail),
			}
			rc, err := client.GetPodLogs(params)
			if err != nil {
				logrus.WithError(err).Errorf("error while getting logs")
				ferr.Println(err)
				ctx.Exit(1)
			}
			defer rc.Close()
			scanner := bufio.NewScanner(rc)
			var nLines uint64
			for scanner.Scan() {
				if err := scanner.Err(); err != nil {
					err = ErrUnableToReadLogs.Wrap(err)
					logrus.WithError(err).Errorf("unable to scan logs byte stream")
					activekit.Attention(err.Error())
					ctx.Exit(1)
				}
				fmt.Println(scanner.Text())
				nLines++
			}
			if err != nil {
				activekit.Attention(err.Error())
				ctx.Exit(1)
			}
		},
		PostRun: ctx.CobraPostRun,
	}
	command.PersistentFlags().
		BoolVarP(&logsConfig.Follow, "follow", "f", false, `follow pod logs`)
	command.PersistentFlags().
		UintVarP(&logsConfig.Tail, "tail", "t", 100, `print last <value> log lines`)
	return command
}
