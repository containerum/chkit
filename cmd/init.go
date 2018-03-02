package cmd

import (
	"os"
	"path"

	"github.com/blang/semver"

	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/model"
	"github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v2"
)

const (
	Version        = "3.0.0-alpha"
	containerumAPI = "https://94.130.09.147:8082"
	FlagConfigFile = "config"
	FlagAPIaddr    = "apiaddr"
)

func Run(args []string) error {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{}
	log.Level = logrus.InfoLevel

	configPath, err := getConfigPath()
	if err != nil {
		log.WithError(err).
			Errorf("error while getting homedir path")
		return err
	}
	var App = &cli.App{
		Name:    "chkit",
		Usage:   "containerum cli",
		Version: semver.MustParse(Version).String(),
		Action: func(ctx *cli.Context) error {
			clientConfig := model.ClientConfig{}
			err := loadConfig(ctx.String("config"), &clientConfig)
			if clientConfig.APIaddr == "" {
				clientConfig.APIaddr = ctx.String("api")
			}
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
			err = saveConfig(configPath, &clientConfig)
			if err != nil {
				log.WithError(err).
					Errorf("error while saving config")
				return err
			}

			return nil
		},
		Metadata: map[string]interface{}{
			"client": client.Client{},
		},
		Commands: []*cli.Command{
			commandLogin(log, configPath),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Usage:   "config file",
				Aliases: []string{"c"},
				Value:   path.Join(configPath, "config.toml"),
			},
			&cli.StringFlag{
				Name:    "api",
				Usage:   "API address",
				Value:   containerumAPI,
				Hidden:  true,
				EnvVars: []string{"CONTAINERUM_API"},
			},
		},
	}
	return App.Run(args)
}
