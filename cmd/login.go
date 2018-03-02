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
		log := getLog(ctx)
		chClient := getClient(ctx)
		if ctx.IsSet("username") {
			chClient.Config.Username = ctx.String("username")
		} else {
			chClient.Config.Username = readLogin(log)
		}
		if ctx.IsSet("pass") {
			chClient.Config.Password = ctx.String("pass")
		} else {
			chClient.Config.Password = readPassword(log)
		}
		setClient(ctx, chClient)
		saveConfig(ctx, chClient.Config)
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
