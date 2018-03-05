package cmd

import (
	"os"
	"path"

	kubeClientModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/blang/semver"
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/model"
	"github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v2"
)

const (
	Version        = "3.0.0-alpha"
	containerumAPI = "http://192.168.88.200" //"https://94.130.09.147:8082"
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
			log := getLog(ctx)
			if err := setupConfig(ctx); err != nil && !os.IsNotExist(err) {
				log.Fatalf("error while config setup: %v", err)
			} else if os.IsNotExist(err) {
				if err := login(ctx); err != nil {
					log.Fatalf("error while login: %v", err)
				}
				config := getConfig(ctx)
				if config.APIaddr == "" {
					config.APIaddr = ctx.String("api")
				}
			}
			if err := setupClient(ctx); err != nil {
				log.Fatalf("error while client setup: %v", err)
			}
			if err := persist(ctx); err != nil {
				log.Fatalf("%v", err)
			}
			clientConfig := getClient(ctx).Config
			log.Infof("logged as %q", clientConfig.Username)
			if err := mainActivity(ctx); err != nil {
				log.Fatalf("error in main activity: %v", err)
			}
			return nil
		},
		After: func(ctx *cli.Context) error {
			return nil
		},
		Metadata: map[string]interface{}{
			"client":     chClient.Client{},
			"configPath": configPath,
			"log":        log,
			"config":     model.ClientConfig{},
			"tokens":     kubeClientModels.Tokens{},
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
