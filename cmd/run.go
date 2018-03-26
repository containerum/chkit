package cmd

import (
	"path"
	"time"

	"github.com/containerum/chkit/cmd/config_dir"

	kubeClientModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/blang/semver"
	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/model"
	"github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v2"
)

var (
	Version        = "3.0.0-alpha"
	FlagConfigFile = "config"
	FlagAPIaddr    = "apiaddr"
)

var (
	ErrFatalError chkitErrors.Err = "fatal error"
)

func Run(args []string) error {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC1123,
	}
	var App = &cli.App{
		Name:    "chkit",
		Usage:   "containerum cli",
		Version: semver.MustParse(Version).String(),
		Action:  runAction,
		Metadata: map[string]interface{}{
			"client":     chClient.Client{},
			"configPath": confDir.ConfigDir(),
			"log":        log,
			"config":     model.Config{},
			"tokens":     kubeClientModels.Tokens{},
			"version":    semver.MustParse(Version),
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
				Value:   path.Join(confDir.ConfigDir(), "config.toml"),
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
				Value:  "",
				Hidden: true,
			},

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
	return App.Run(args)
}

func runAction(ctx *cli.Context) error {
	if err := setupLog(ctx); err != nil {
		return err
	}
	log := util.GetLog(ctx)
	log.Debugf("loading config")
	if err := loadConfig(ctx); err != nil {
		return err
	}
	log.Debugf("running setup")
	err := setupConfig(ctx)
	config := util.GetConfig(ctx)
	switch {
	case err == nil:
		// pass
	case ErrInvalidUserInfo.Match(err):
		log.Debugf("invalid user information")
		log.Debugf("running login")
		user, err := login(ctx)
		if err != nil {
			return err
		}
		config.UserInfo = user
		util.SetConfig(ctx, config)
	default:
		log.Debugf("fatal error")
		return ErrFatalError.Wrap(err)
	}
	log.Debugf("client initialisation")
	if err := setupClient(ctx); err != nil {
		return err
	}
	if err := persist(ctx); err != nil {
		log.Fatalf("%v", err)
	}
	client := util.GetClient(ctx)
	if err := util.SaveTokens(ctx, client.Tokens); err != nil {
		return chkitErrors.NewExitCoder(err)
	}
	clientConfig := client.Config
	log.Infof("Hello, %q!", clientConfig.Username)
	if err := mainActivity(ctx); err != nil {
		log.Fatalf("error in main activity: %v", err)
	}
	return nil
}
