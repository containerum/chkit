package cmd

import (
	"os"

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
	log := getLog(ctx)
	config := getConfig(ctx)
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
	setClient(ctx, *client)
	return nil
}

func setupConfig(ctx *cli.Context) error {
	log := getLog(ctx)
	config := getConfig(ctx)
	err := loadConfig(ctx.String("config"), &config)
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
	setConfig(ctx, config)
	return nil
}

func persist(ctx *cli.Context) error {
	if !ctx.IsSet("config") {
		return saveConfig(ctx)
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
	client := getClient(ctx)
	tokens, err := loadTokens(ctx)
	if err != nil {
		return err
	}
	client.Tokens = tokens
	if err := client.Auth(); err != nil {
		return err
	}
	return saveTokens(ctx, tokens)
}
