package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/configuration"
	. "github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/util/activeToolkit"
	"github.com/containerum/chkit/pkg/util/animation"
	"github.com/containerum/chkit/pkg/util/trasher"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
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
)

var commandLogin = &cobra.Command{
	Use: "login",
	Run: func(command *cobra.Command, args []string) {
		err := setupConfig()
		switch {
		case err == nil || ErrInvalidUserInfo.Match(err) || ErrUnableToLoadTokens.Match(err):
			userInfo, err := login()
			if err != nil {
				logrus.WithError(err).Errorf("unable to setup config")
				fmt.Printf("Unable to setup config :(\n")
				return
			}
			Context.Client.UserInfo = userInfo
		default:
			panic(ErrFatalError.Wrap(err))
		}
		if err := setupClient(); err != nil {
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
		mainActivity()
	},
}

func login() (model.UserInfo, error) {
	user := model.UserInfo{}
	var err error
	if strings.TrimSpace(runContext.Username) != "" {
		user.Username = runContext.Username
	} else {
		user.Username, err = readLogin()
		if err != nil {
			return user, err
		}
	}
	if strings.TrimSpace(user.Username) == "" {
		return user, ErrInvalidUsername
	}

	if strings.TrimSpace(runContext.Pass) != "" {
		user.Password = runContext.Pass
	} else {
		user.Password, err = readPassword()
		if err != nil {
			return user, err
		}
	}
	if strings.TrimSpace(user.Password) == "" {
		return user, ErrInvalidPassword
	}
	return user, nil
}

func readLogin() (string, error) {
	fmt.Print("Enter your email: ")
	email, err := bufio.NewReader(os.Stdin).ReadString('\n')
	email = strings.TrimRight(email, "\r\n")
	if err != nil {
		return "", ErrUnableToReadUsername.Wrap(err)
	}
	return email, nil
}

func readPassword() (string, error) {
	fmt.Print("Enter your password: ")
	passwordB, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", ErrUnableToReadPassword.Wrap(err)
	}
	fmt.Println("")
	return string(passwordB), nil
}
