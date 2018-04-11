package login

import (
	"fmt"
	"time"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/cli/clisetup"
	"github.com/containerum/chkit/pkg/configuration"
	. "github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activeToolkit"
	"github.com/containerum/chkit/pkg/util/animation"
	"github.com/containerum/chkit/pkg/util/trasher"
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

var Command = &cobra.Command{
	Use: "login",
	Run: func(command *cobra.Command, args []string) {
		err := clisetup.SetupConfig()
		switch {
		case err == nil || clisetup.ErrInvalidUserInfo.Match(err) || clisetup.ErrUnableToLoadTokens.Match(err):
			err := Login()
			if err != nil {
				logrus.WithError(err).Errorf("unable to setup config")
				fmt.Printf("Unable to setup config :(\n")
				return
			}
			Context.Changed = true
		default:
			panic(ErrFatalError.Wrap(err))
		}
		if err := clisetup.SetupClient(); err != nil {
			logrus.WithError(err).Errorf("unable to setup client")
			panic(err)
		}
		anim := &animation.Animation{
			Framerate:      0.5,
			Source:         trasher.NewSilly(),
			ClearLastFrame: true,
		}
		go func() {
			time.Sleep(4 * time.Second)
			anim.Run()
		}()
		err = func() error {
			defer anim.Stop()
			return Context.Client.Auth()
		}()
		if err != nil {
			logrus.WithError(err).Errorf("unable to auth")
			fmt.Printf("Unable to authenticate :(")
			return
		}
		if err := configuration.SaveTokens(Context.Client.Tokens); err != nil {
			logrus.WithError(err).Errorf("unable to save tokens")
			fmt.Printf("Unable to save tokens :(")
		}
		Context.Namespace, err = configuration.GetFirstClientNamespace()

		if err != nil {
			logrus.WithError(err).Error("unable to get default namespace")
			if !Context.Quiet {
				fmt.Printf("Unable to get default namespace :(\n")
				_, option, _ := activeToolkit.Options("What's next?", false,
					"Choose your own namespace",
					"Exit")
				switch option {
				case 0:
					Context.Namespace, _ = activeToolkit.AskLine("Type namespace: ")
				default:
					// pass
				}
			}
		}
		if !Context.Quiet {
			fmt.Printf("Successfuly authenticated as %q ^_^\n", Context.Client.Username)
			fmt.Printf("Using %q as default namespace\n", Context.Namespace)
		}
	},
	PostRun: func(command *cobra.Command, args []string) {
		if Context.Changed {
			if err := configuration.SaveConfig(); err != nil {
				logrus.WithError(err).Errorf("unable to save config")
				fmt.Printf("Unable to save config: %v\n", err)
				return
			}
		}
		if err := configuration.SaveTokens(Context.Client.Tokens); err != nil {
			logrus.WithError(err).Errorf("unable to save tokens")
			fmt.Printf("Unable to save tokens: %v\n", err)
			return
		}
	},
}
