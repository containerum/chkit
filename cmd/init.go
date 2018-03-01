package cmd

import (
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

var (
	Configuration = model.Config{}
	log           = &logrus.Logger{
		Formatter: &logrus.TextFormatter{},
	}
)

func Run(args []string) error {

	var App = &cli.App{
		Name:    "chkit",
		Version: semver.MustParse(Version).String(),
		Action: func(ctx *cli.Context) error {

			return nil
		},
		Before: func(ctx *cli.Context) error {
			err := initConfig()
			if err != nil {
				log.WithError(err).
					Errorf("error while getting homedir path")
				return err
			}
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   Configuration.ConfigPath,
			},
		},
	}
	return App.Run(args)
}
