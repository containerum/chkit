package cmd

import (
	"os"

	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/model"
	"github.com/labstack/gommon/log"
	"gopkg.in/urfave/cli.v2"
)

const (
	ErrUnableToLoadConfig chkitErrors.Err = "unable to load config"
	ErrInvalidUserInfo    chkitErrors.Err = "invalid user info"
	ErrInvalidAPIurl      chkitErrors.Err = "invalid API url"
	ErrUnableToLoadTokens chkitErrors.Err = "unable to load tokens"
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
		log.Infof("Using test API")
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
	if ctx.IsSet("test") {
		testAPIurl := os.Getenv("CONTAINERUM_API")
		log.Infof("using test api %q", testAPIurl)
		config.APIaddr = testAPIurl
	}
	if config.APIaddr == "" {
		log.Debug("invalid API url")
		return ErrInvalidAPIurl
	}
	if config.Password == "" || config.Username == "" {
		log.Debugf("invalid username or pass")
		return ErrInvalidUserInfo
	}
	config.Fingerprint = Fingerprint()
	tokens, err := util.LoadTokens(ctx)
	if err != nil && !os.IsNotExist(err) {
		return ErrUnableToLoadTokens.Wrap(err)
	} else if os.IsNotExist(err) {
		err = util.SaveTokens(ctx, model.Tokens{})
		if err != nil {
			return err
		}
	}
	config.Tokens = tokens
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

	err := os.MkdirAll(util.GetConfigPath(ctx), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	_, err = os.Stat(ctx.String("config"))
	if err != nil && os.IsNotExist(err) {
		file, err := os.Create(ctx.String("config"))
		if err != nil {
			return err
		}
		return file.Close()
	} else if err != nil {
		return err
	}

	err = util.LoadConfig(ctx.String("config"), &config)
	if err != nil {
		return err
	}
	return nil
}

func setupAll(ctx *cli.Context) error {
	log := util.GetLog(ctx)
	log.Debugf("setuping config")
	if err := setupConfig(ctx); err != nil {
		return err
	}
	log.Debugf("setuping client")
	if err := setupClient(ctx); err != nil {
		return err
	}
	return nil
}
