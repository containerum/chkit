package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"

	"golang.org/x/crypto/ssh/terminal"
	cli "gopkg.in/urfave/cli.v2"
)

var commandLogin = &cli.Command{
	Name:  "login",
	Usage: "use username and password to login in the system",
	Action: func(ctx *cli.Context) error {
		login(ctx)
		if err := setupClient(ctx); err != nil {
			return err
		}
		return nil
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "pass",
			Usage: "password to system",
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "your account email",
		},
	},
}

func login(ctx *cli.Context) {
	log := getLog(ctx)
	config := getConfig(ctx)
	if ctx.IsSet("username") {
		config.Username = ctx.String("username")
	} else {
		config.Username = readLogin(log)
	}
	if ctx.IsSet("pass") {
		config.Password = ctx.String("pass")
	} else {
		config.Password = readPassword(log)
	}
	setConfig(ctx, config)
}

func readLogin(log *logrus.Logger) string {
	fmt.Print("Enter your email: ")
	email, err := bufio.NewReader(os.Stdin).ReadString('\n')
	email = strings.TrimRight(email, "\r\n")
	exitOnErr(log, err)
	return email
}

func readPassword(log *logrus.Logger) string {
	fmt.Print("Enter your password: ")
	passwordB, err := terminal.ReadPassword(int(syscall.Stdin))
	exitOnErr(log, err)
	return string(passwordB)
}
