package cmd

import (
	"path"

	kubeClientModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/blang/semver"
	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/model"
	"github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v2"
)

const (
	Version        = "3.0.0-alpha"
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
			log := util.GetLog(ctx)
			switch err := setupAll(ctx).(type) {
			case nil:
			default:
				return err
			}
			if err := setupClient(ctx); err != nil {
				log.Fatalf("error while client setup: %v", err)
			}
			tokens, err := util.LoadTokens(ctx)
			if err != nil {
				return chkitErrors.NewExitCoder(err)
			}
			client := util.GetClient(ctx)
			client.Tokens = tokens
			if err := client.Auth(); err != nil {
				return err
			}
			if err := persist(ctx); err != nil {
				log.Fatalf("%v", err)
			}
			if err := util.SaveTokens(ctx, tokens); err != nil {
				return chkitErrors.NewExitCoder(err)
			}
			clientConfig := client.Config
			log.Infof("Hello, %q!", clientConfig.Username)
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
			"config":     model.Config{},
			"tokens":     kubeClientModels.Tokens{},
		},
		Commands: []*cli.Command{
			commandLogin,
			commandGet,
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
				Value:   "",
				Hidden:  true,
				EnvVars: []string{"CONTAINERUM_API"},
			},
			&cli.StringFlag{
				Name:   "test",
				Usage:  "test presets",
				Value:  "api",
				Hidden: true,
			},
		},
	}
	return App.Run(args)
}
