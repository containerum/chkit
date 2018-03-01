package cmd

import (
	"os"
	"path"

	"github.com/blang/semver"

	"github.com/containerum/chkit/pkg/model"
	"github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v2"
)

const (
	Version        = "3.0.0-alpha"
	FlagAPIaddr    = "apiaddr"
	FlagConfigFile = "config"
)

func Run(args []string) error {
	Configuration := &model.Config{}
	log := &logrus.Logger{
		Formatter: &logrus.TextFormatter{},
	}
	err := initConfig(Configuration)
	if err != nil {
		log.WithError(err).
			Errorf("error while getting homedir path")
		return err
	}
	var App = &cli.App{
		Name:    "chkit",
		Version: semver.MustParse(Version).String(),
		Action: func(ctx *cli.Context) error {
			err := loadConfig(&Configuration.Client, ctx.String("config"))
			if err != nil && !os.IsNotExist(err) {
				log.WithError(err).
					Errorf("error while loading config file")
				return err
			} else if os.IsNotExist(err) {
				log.Info("You are not logged in!")
				err = ctx.App.Command("login").Run(ctx)
				if err != nil {
					return err
				}
			}
			return nil
		},
		Commands: []*cli.Command{
			commandLogin(log, Configuration),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Usage:   "config file",
				Aliases: []string{"c"},
				Value:   path.Join(Configuration.ConfigPath, "config.toml"),
			},
		},
	}
	return App.Run(args)
}
