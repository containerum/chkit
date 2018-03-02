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

func login(ctx *cli.Context) {
	log := getLog(ctx)
	config := getConfig(ctx)
	if ctx.IsSet("username") {
		config.Username = ctx.String("username")
	} else {
		config.Username = readLogin(log)
	}
	if strings.TrimSpace(config.Username) == "" {
		log.Fatalln("Username must be non-empty string!")
	}

	if ctx.IsSet("pass") {
		config.Password = ctx.String("pass")
	} else {
		config.Password = readPassword(log)
	}
	if strings.TrimSpace(config.Password) == "" {
		log.Fatalln("Password must be non-empty string!")
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
