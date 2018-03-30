package cmd

import (
	"fmt"
	"os"
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
	// Version -- chkit version
	Version = "3.0.0-alpha"
)

const (
	// FlagConfigFile -- context config data key
	FlagConfigFile = "config"
	// FlagAPIaddr -- API address context key
	FlagAPIaddr = "apiaddr"
)

var (
	// ErrFatalError -- unrecoverable fatal error
	ErrFatalError chkitErrors.Err = "fatal error"
)

// Run -- root action
func Run(args []string) error {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC1123,
	})

	if !DEBUG {
		logFile := path.Join(confDir.ConfigDir(), util.LogFileName())
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			logrus.Fatalf("error while creating log file: %v", err)
		}
		logrus.SetOutput(file)
	}
	var App = &cli.App{
		Name:    "chkit",
		Usage:   "containerum cli",
		Version: semver.MustParse(Version).String(),
		/*Before: func(ctx *cli.Context) error {
			var updater update.LatestCheckerDownloader
			currVersion := semver.MustParse(Version)
			updater = update.NewGithubLatestCheckerDownloader(ctx, "containerum", "chkit")
			version, err := updater.LatestVersion()
			if err != nil {
				return err
			}
			if currVersion.LE(version) {
				if yes, err := update.AskForUpdate(ctx, version); err != nil {
					return err
				} else if yes {
					if err := update.Update(ctx, updater, true); err != nil {
						return err
					}
				}
			}
			return nil
		},*/
		Action: runAction,
		Metadata: map[string]interface{}{
			"client":     chClient.Client{},
			"configPath": confDir.ConfigDir(),
			"config":     model.Config{},
			"tokens":     kubeClientModels.Tokens{},
			"version":    semver.MustParse(Version),
		},
		Commands: []*cli.Command{
			&cli.Command{
				Name:        "version",
				Usage:       "prints chkit version",
				Description: "prints chkit version. Aliases: vers, vs, v",
				Aliases:     []string{"vers", "vs", "v"},
				Action: func(ctx *cli.Context) error {
					fmt.Println(Version)
					return nil
				},
			},
			commandLogin,
			commandGet,
			commandDelete,
			commandUpdate,
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
			&cli.BoolFlag{
				Name:   "debug-requests",
				Value:  false,
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
			&cli.StringFlag{
				Name:    "namespace",
				Aliases: []string{"n"},
				Usage:   "namespace to use",
			},
		},
	}
	return App.Run(args)
}

func runAction(ctx *cli.Context) error {
	logrus.Debugf("loading config")
	if err := loadConfig(ctx); err != nil {
		return err
	}
	logrus.Debugf("running setup")
	err := setupConfig(ctx)
	config := util.GetConfig(ctx)
	switch {
	case err == nil:
		// pass
	case ErrInvalidUserInfo.Match(err):
		logrus.Debugf("invalid user information")
		logrus.Debugf("running login")
		user, err := login(ctx)
		if err != nil {
			return err
		}
		config.UserInfo = user
		util.SetConfig(ctx, config)
	default:
		logrus.Debugf("fatal error")
		return ErrFatalError.Wrap(err)
	}
	logrus.Debugf("client initialisation")
	if err := setupClient(ctx); err != nil {
		return err
	}
	client := util.GetClient(ctx)
	if err := util.SaveTokens(ctx, client.Tokens); err != nil {
		return chkitErrors.NewExitCoder(err)
	}
	config.DefaultNamespace, err = util.GetFirstClientNamespace(ctx)
	if err != nil {
		return err
	}
	util.SetConfig(ctx, config)
	if err := persist(ctx); err != nil {
		logrus.Fatalf("%v", err)
	}
	// re-setup client to save default namespace
	if err := setupClient(ctx); err != nil {
		return err
	}
	clientConfig := client.Config
	logrus.Infof("Hello, %q!", clientConfig.Username)
	if err := mainActivity(ctx); err != nil {
		logrus.Fatalf("error in main activity: %v", err)
	}
	return nil
}
