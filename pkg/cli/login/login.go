package login

import (
	"os"

	"fmt"

	"github.com/containerum/chkit/pkg/cli/clisetup"
	"github.com/containerum/chkit/pkg/cli/postrun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/spf13/cobra"
)

func Login(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:   "login",
		Short: "Login to system",
		Run: func(command *cobra.Command, args []string) {
			if err := clisetup.SetupLogs(ctx); err != nil {
				angel.Angel(ctx, err)
				os.Exit(1)
			}
			flags := command.Flags()
			if flags.Changed("default-namespace") {
				defNS, _ := flags.GetString("default-namespace")
				ns, err := ctx.Client.GetNamespace(defNS)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				ctx.SetNamespace(ns)
			}
			if err := Setup(ctx); err != nil {
				angel.Angel(ctx, err)
				os.Exit(1)
			}
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			postrun.PostRun(coblog.Logger(cmd), ctx)
		},
	}
	command.PersistentFlags().
		StringVarP(&ctx.Client.Username, "username", "u", "", "your account login")
	command.PersistentFlags().
		StringVarP(&ctx.Client.Password, "password", "p", "", "your account password")
	command.PersistentFlags().
		String("default-namespace", "", "use as default namespace, if '-', then use first one")
	return command
}
