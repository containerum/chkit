package cmd

import (
	"os"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/client"
	"gopkg.in/urfave/cli.v2"
)

const (
	ErrUnableToLoadConfig chkitErrors.Err = "unable to load config"
	ErrInvalidAPIurl      chkitErrors.Err = "invalid API url"
	ErrUnableToInitClient chkitErrors.Err = "unable to init client"
)

func setupClient(ctx *cli.Context) error {
	log := getLog(ctx)
	config := getConfig(ctx)
	var client *chClient.Client
	var err error
	if ctx.Bool("test") {
		log.Infof("running in test mode")
		client, err = chClient.NewClient(config, chClient.UnsafeSkipTLSCheck)
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
	config := getConfig(ctx)
	err := loadConfig(ctx.String("config"), &config)
	if err != nil {
		return ErrUnableToLoadConfig.
			Wrap(err)
	}
	if ctx.Bool("test") {
		config.APIaddr = os.Getenv("CONTAINERUM_API")
	}
	if config.APIaddr == "" {
		return ErrInvalidAPIurl
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
