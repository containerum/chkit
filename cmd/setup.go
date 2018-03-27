package cmd

import (
	"net/url"
	"os"

	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/model"
	"github.com/ninedraft/delog"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

const (
	ErrUnableToLoadConfig chkitErrors.Err = "unable to load config"
	ErrInvalidUserInfo    chkitErrors.Err = "invalid user info"
	ErrInvalidAPIurl      chkitErrors.Err = "invalid API url"
	ErrUnableToLoadTokens chkitErrors.Err = "unable to load tokens"
	ErrUnableToSaveTokens chkitErrors.Err = "unable to save tokens"
)

func setupClient(ctx *cli.Context) error {
	log := util.GetLog(ctx)
	config := util.GetConfig(ctx)
	var client *chClient.Client
	var err error
	switch ctx.String("test") {
	case "mock":
		log.Infof("Using mock API")
		client, err = chClient.NewClient(config, chClient.WithMock)
	case "api":
		log.Debugf("Using test API: %q", config.APIaddr)
		config.Log = log.WriterLevel(logrus.DebugLevel)
		client, err = chClient.NewClient(config, chClient.WithTestAPI)
	default:
		client, err = chClient.NewClient(config, chClient.WithCommonAPI)
	}
	if err != nil {
		return err
	}
	util.SetClient(ctx, *client)
	return nil
}

func setupConfig(ctx *cli.Context) error {
	config := util.GetConfig(ctx)
	log := util.GetLog(ctx)
	log.Debugf("test: %q", ctx.String("test"))
	config.Fingerprint = Fingerprint()
	tokens, err := util.LoadTokens(ctx)
	if err != nil && !os.IsNotExist(err) {
		return ErrUnableToLoadTokens.Wrap(err)
	} else if os.IsNotExist(err) {
		if err = util.SaveTokens(ctx, model.Tokens{}); err != nil {
			return ErrUnableToSaveTokens.Wrap(err)
		}
	}
	config.Tokens = tokens
	if ctx.IsSet("test") {
		testAPIurl := os.Getenv("CONTAINERUM_API")
		log.Debugf("using test api %q", testAPIurl)
		config.APIaddr = testAPIurl
	}
	if _, err := url.Parse(config.APIaddr); err != nil {
		log.Debugf("invalid API url: %q", config.APIaddr)
		return ErrInvalidAPIurl.Wrap(err)
	}
	if config.Password == "" || config.Username == "" {
		log.Debugf("invalid username or pass")
		util.SetConfig(ctx, config)
		return ErrInvalidUserInfo
	}
	util.SetConfig(ctx, config)
	return nil
}

func persist(ctx *cli.Context) error {
	if !ctx.IsSet("config") {
		return util.SaveConfig(ctx)
	}
	return nil
}

func loadConfig(ctx *cli.Context) error {
	//log := util.GetLog(ctx)
	config := util.GetConfig(ctx)
	err := util.LoadConfig(ctx.String("config"), &config)
	if err != nil {
		return ErrUnableToLoadConfig.Wrap(err)
	}
	util.SetConfig(ctx, config)
	return nil
}

func setupAll(ctx *cli.Context) error {
	log := util.GetLog(ctx)
	log.Debugf("setuping config")
	if err := loadConfig(ctx); err != nil {
		return err
	}
	if err := setupConfig(ctx); err != nil {
		return err
	}
	log.Debugf("setuping client")
	if err := setupClient(ctx); err != nil {
		return err
	}
	client := util.GetClient(ctx)
	log.Debugf("API: %q", client.APIaddr)
	return nil
}

func setupLog(ctx *cli.Context) error {
	log := util.GetLog(ctx)
	if ctx.IsSet("test") {
		log.Formatter = delog.NewFormatter(log.Formatter)
		log.SetLevel(logrus.DebugLevel)
		log.Debug("debug mode on")
	}
	return nil
}
