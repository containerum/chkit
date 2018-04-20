package login

import (
	"fmt"

	"os"

	"github.com/containerum/chkit/pkg/cli/clisetup"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Login(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use: "login",
		Run: func(command *cobra.Command, args []string) {
			if err := clisetup.Setup(ctx); err != nil {
				angel.Angel(ctx, err)
				os.Exit(1)
			}
		},
		PostRun: func(command *cobra.Command, args []string) {
			if ctx.Changed {
				if err := configuration.SaveConfig(ctx); err != nil {
					logrus.WithError(err).Errorf("unable to save config")
					fmt.Printf("Unable to save config: %v\n", err)
					return
				}
			}
			if err := configuration.SaveTokens(ctx, ctx.Client.Tokens); err != nil {
				logrus.WithError(err).Errorf("unable to save tokens")
				fmt.Printf("Unable to save tokens: %v\n", err)
				return
			}
		},
	}
	command.PersistentFlags().
		StringVarP(&ctx.Client.Username, "username", "u", "", "your account login")
	command.PersistentFlags().
		StringVarP(&ctx.Client.Password, "password", "p", "", "your account password")

	return command
}
