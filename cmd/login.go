package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/containerum/chkit/pkg/util/animation"
	"github.com/containerum/chkit/pkg/util/trasher"

	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"
	"github.com/sirupsen/logrus"

	"golang.org/x/crypto/ssh/terminal"
	cli "gopkg.in/urfave/cli.v2"
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

var commandLogin = &cli.Command{
	Name:  "login",
	Usage: "login your in the system",
	Action: func(ctx *cli.Context) error {
		err := setupConfig(ctx)
		config := util.GetConfig(ctx)
		switch {
		case err == nil || ErrInvalidUserInfo.Match(err) || ErrUnableToLoadTokens.Match(err):
			userInfo, err := login(ctx)
			if err != nil {
				logrus.Debugf("fatal error: %v", err)
				return err
			}
			config.UserInfo = userInfo
			util.SetConfig(ctx, config)
		default:
			return err
		}
		if err := setupClient(ctx); err != nil {
			return err
		}
		client := util.GetClient(ctx)
		client.Tokens = model.Tokens{}

		anim := &animation.Animation{
			Framerate:      0.5,
			Source:         trasher.NewSilly(),
			ClearLastFrame: true,
		}
		go anim.Run()
		if err := client.Auth(); err != nil {
			anim.Stop()
			fmt.Println(err)
			return err
		}
		anim.Stop()
		if err := util.SaveTokens(ctx, client.Tokens); err != nil {
			return err
		}
		config.DefaultNamespace, err = util.GetFirstClientNamespace(ctx)
		if err != nil {
			return err
		}
		util.SetConfig(ctx, config)
		if err := persist(ctx); err != nil {
			return err
		}
		return mainActivity(ctx)
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "username",
			Usage: "your account email",
		},
		&cli.StringFlag{
			Name:  "pass",
			Usage: "password to system",
		},
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
