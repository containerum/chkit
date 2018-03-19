package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"

	"golang.org/x/crypto/ssh/terminal"
	cli "gopkg.in/urfave/cli.v2"
)

var (
	ErrUnableToReadPassword chkitErrors.Err = "unable to read password"
	ErrUnableToReadUsername chkitErrors.Err = "unable to read username"
	ErrInvalidPassword      chkitErrors.Err = "invalid password"
	ErrInvalidUsername      chkitErrors.Err = "invalid username"
)

var commandLogin = &cli.Command{
	Name:  "login",
	Usage: "login your in the system",
	Action: func(ctx *cli.Context) error {
		if err := setupConfig(ctx); err != nil {
			return err
		}
		if err := setupClient(ctx); err != nil {
			return err
		}

		client := util.GetClient(ctx)
		client.Tokens = model.Tokens{}
		loginData, err := login(ctx)
		if err != nil {
			return err
		}
		client.Config.UserInfo = loginData
		util.SetClient(ctx, client)

		client = util.GetClient(ctx)
		if err := client.Auth(); err != nil {
			return chkitErrors.NewExitCoder(err)
		}
		if err := util.SaveTokens(ctx, client.Tokens); err != nil {
			return chkitErrors.NewExitCoder(err)
		}
		return mainActivity(ctx)
	},
}

func login(ctx *cli.Context) (model.UserInfo, error) {
	user := model.UserInfo{}
	var err error
	if ctx.IsSet("username") {
		user.Username = ctx.String("username")
	} else {
		user.Username, err = readLogin()
		if err != nil {
			return user, err
		}
	}
	if strings.TrimSpace(user.Username) == "" {
		return user, ErrInvalidUsername
	}

	if ctx.IsSet("pass") {
		user.Password = ctx.String("pass")
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
