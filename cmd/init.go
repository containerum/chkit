package cmd

import (
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

	configPath, err := configPath()
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
			err := setupConfig(ctx)
			if err != nil {
				log.Error(err)
				return err
			}
			clientConfig := getClient(ctx).Config
			log.Infof("logged as %q", clientConfig.Username)
			return err
		},
		Metadata: map[string]interface{}{
			"client":     chClient.Client{},
			"configPath": configPath,
			"log":        log,
			"config":     model.ClientConfig{},
		},
		Commands: []*cli.Command{
			commandLogin,
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
