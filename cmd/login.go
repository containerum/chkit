package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/containerum/chkit/pkg/chkitErrors"

	"golang.org/x/crypto/ssh/terminal"
	cli "gopkg.in/urfave/cli.v2"
)

var commandLogin = &cli.Command{
	Name:  "login",
	Usage: "login your in the system",
	Action: func(ctx *cli.Context) error {
		login(ctx)
		config := getConfig(ctx)
		if config.APIaddr == "" {
			config.APIaddr = ctx.String("api")
		}
		setConfig(ctx, config)
		persist(ctx)
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

func login(ctx *cli.Context) error {
	log := getLog(ctx)
	config := getConfig(ctx)
	var err error
	if ctx.IsSet("username") {
		config.Username = ctx.String("username")
	} else {
		config.Username, err = readLogin()
		if err != nil {
			return err
		}
	}
	if strings.TrimSpace(config.Username) == "" {
		return chkitErrors.ErrInvalidUsername().
			AddDetailF("username must be non-empty string!")
	}

	if ctx.IsSet("pass") {
		config.Password = ctx.String("pass")
	} else {
		config.Password, err = readPassword()
		if err != nil {
			return err
		}
	}
	if strings.TrimSpace(config.Password) == "" {
		return chkitErrors.ErrInvalidPassword().
			AddDetailF("Password must be non-empty string!")
	}
	setConfig(ctx, config)
	return nil
}

func readLogin() (string, error) {
	fmt.Print("Enter your email: ")
	email, err := bufio.NewReader(os.Stdin).ReadString('\n')
	email = strings.TrimRight(email, "\r\n")
	if err != nil {
		return "", chkitErrors.ErrUnableToReadUsername().
			AddDetailsErr(err)
	}
	return email, nil
}

func readPassword() (string, error) {
	fmt.Print("Enter your password: ")
	passwordB, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", chkitErrors.ErrUnableToReadPassword().
			AddDetailsErr(err)
	}
	return string(passwordB), nil
}
