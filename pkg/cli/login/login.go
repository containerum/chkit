package login

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/cli/clisetup"
	"github.com/containerum/chkit/pkg/configdir"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// ErrUnableToReadPassword -- unable to read password
	ErrUnableToReadPassword chkitErrors.Err = "unable to read password"
	// ErrUnableToReadUsername -- unable to read username
	ErrUnableToReadUsername chkitErrors.Err = "unable to read username"
	// ErrInvalidPassword -- invalid password
	ErrInvalidPassword chkitErrors.Err = "invalid password"
	// ErrInvalidUsername -- invalid username
	ErrInvalidUsername chkitErrors.Err = "invalid username"
	// ErrFatalError -- unrecoverable fatal error
	ErrFatalError chkitErrors.Err = "fatal error"
)

func Login(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use: "login",
		PreRun: func(cmd *cobra.Command, args []string) {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.SetFormatter(&logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: time.RFC1123,
			})
			logFile := path.Join(configdir.LogDir(), configuration.LogFileName())
			file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
			if err != nil {
				logrus.Fatalf("error while creating log file: %v", err)
			}
			logrus.SetOutput(file)
		},
		Run: func(command *cobra.Command, args []string) {
			err := clisetup.SetupConfig(ctx)
			switch {
			case err == nil || clisetup.ErrInvalidUserInfo.Match(err) || clisetup.ErrUnableToLoadTokens.Match(err):
				err := InteractiveLogin(ctx)
				if err != nil {
					logrus.WithError(err).Errorf("unable to setup config")
					fmt.Printf("Unable to setup config :(\n")
					return
				}
				ctx.Changed = true
			default:
				panic(ErrFatalError.Wrap(err))
			}
			if err := clisetup.SetupClient(ctx, false); err != nil {
				logrus.WithError(err).Errorf("unable to setup client")
				panic(err)
			}
			err = func() error {
				return ctx.Client.Auth()
			}()
			if err != nil {
				logrus.WithError(err).Errorf("unable to auth")
				fmt.Printf("Unable to authenticate :(")
				return
			}
			if err := configuration.SaveTokens(ctx, ctx.Client.Tokens); err != nil {
				logrus.WithError(err).Errorf("unable to save tokens")
				fmt.Printf("Unable to save tokens :(")
			}
			ctx.Namespace, err = configuration.GetFirstClientNamespace(ctx)

			if err != nil {
				logrus.WithError(err).Error("unable to get default namespace")
				if !ctx.Quiet {
					fmt.Printf("Unable to get default namespace :(\n")
					(&activekit.Menu{
						Items: []*activekit.MenuItem{
							{
								Label: "Choose your own namespace",
								Action: func() error {
									ctx.Namespace = activekit.Promt("Type namespace label: ")
									return nil
								},
							},
							{
								Label: "Exit",
							},
						},
					}).Run()
				}
			}
			if !ctx.Quiet {
				fmt.Printf("Successfuly authenticated as %q ^_^\n", ctx.Client.Username)
				if ctx.Namespace != "" {
					fmt.Printf("Using %q as default namespace\n", ctx.Namespace)
				} else {
					fmt.Printf("Default namespace is not defined\n")
				}
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
