package cmd

import (
	"os"

	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/client"
	"gopkg.in/urfave/cli.v2"
)

const (
	ErrUnableToLoadConfig chkitErrors.Err = "unable to load config"
	ErrInvalidUserInfo    chkitErrors.Err = "invalid user info"
	ErrInvalidAPIurl      chkitErrors.Err = "invalid API url"
	ErrUnableToInitClient chkitErrors.Err = "unable to init client"
)

func setupClient(ctx *cli.Context) error {
	log := util.GetLog(ctx)
	config := util.GetConfig(ctx)
	var client *chClient.Client
	var err error
	if ctx.IsSet("test") {
		log.Infof("running in test mode")
		switch ctx.String("test") {
		case "mock":
			client, err = chClient.NewClient(config, chClient.Mock)
		case "api":
			client, err = chClient.NewClient(config, chClient.UnsafeSkipTLSCheck)
		}
	} else {
		client, err = chClient.NewClient(config)
	}
	if err != nil {
		return ErrUnableToInitClient.
			Wrap(err)
	}
	util.SetClient(ctx, *client)
	return nil
}

func setupConfig(ctx *cli.Context) error {
	log := util.GetLog(ctx)
	config := util.GetConfig(ctx)
	err := util.LoadConfig(ctx.String("config"), &config)
	if err != nil {
		return ErrUnableToLoadConfig.
			Wrap(err)
	}
	if ctx.IsSet("test") {
		testAPIurl := os.Getenv("CONTAINERUM_API")
		log.Infof("using test api %q", testAPIurl)
		config.APIaddr = testAPIurl
	}

	if config.APIaddr == "" {
		return ErrInvalidAPIurl
	}
	if config.Password == "" || config.Username == "" {
		return ErrInvalidUserInfo
	}

	config.Fingerprint = Fingerprint()
	util.SetConfig(ctx, config)
	return nil
}

func persist(ctx *cli.Context) error {
	if !ctx.IsSet("config") {
		return util.SaveConfig(ctx)
	}
	return nil
}

func setupAll(ctx *cli.Context) error {
	if err := setupConfig(ctx); err != nil {
		return err
	}
	if err := setupClient(ctx); err != nil {
		return err
	}
	client := util.GetClient(ctx)
	tokens, err := util.LoadTokens(ctx)
	if err != nil {
		return err
	}
	client.Tokens = tokens
	util.SetClient(ctx, client)
	return nil
}
