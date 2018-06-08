package login

import (
	"os"

	"fmt"

	"strings"

	"github.com/containerum/chkit/pkg/cli/clisetup"
	"github.com/containerum/chkit/pkg/cli/postrun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Login(ctx *context.Context) *cobra.Command {
	var flags struct {
		Username string `flag:"username u"`
		Password string `flag:"password p"`
	}
	command := &cobra.Command{
		Use:   "login",
		Short: "Login to system",
		Run: func(command *cobra.Command, args []string) {
			if err := clisetup.SetupLogs(ctx); err != nil {
				angel.Angel(ctx, err)
				os.Exit(1)
			}
			var logger = coblog.Logger(command)
			logger.Debugf("start")
			defer logger.Debugf("end")
			ctx.Client.Username = flags.Username
			ctx.Client.Password = flags.Password
			logger.Debugf("start app setup")
			if err := Setup(ctx); err != nil {
				angel.Angel(ctx, err)
				os.Exit(1)
			}
			logger.Debugf("end setup")
			var ns, _ = command.Flags().GetString("namespace")
			switch ns {
			case "-":
				clisetup.GetDefaultNS(ctx, true)
			case "":
				clisetup.GetDefaultNS(ctx, false)
			default:
				nsList, err := ctx.Client.GetNamespaceList()
				logger.Debugf("Getting namespace list")
				if err != nil {
					logger.WithError(err).Errorf("unable to get namespace lsit")
					fmt.Println(err)
					os.Exit(1)
				}
				_, ok := nsList.GetByUserFriendlyID(ns)
				if !ok {
					fmt.Printf("Namespace %q not found!\n%s",
						ns, strings.Join(nsList.OwnersAndLabels(), "\n"))
					os.Exit(1)
				}
			}
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			postrun.PostRun(coblog.Logger(cmd), ctx)
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
