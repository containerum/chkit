package cmd

import (
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/client"

	"gopkg.in/urfave/cli.v2"
)

func setupClient(ctx *cli.Context) error {
	log := getLog(ctx)
	config := getConfig(ctx)
	if config.APIaddr == "" {
		config.APIaddr = ctx.String("api")
	}
	client, err := chClient.NewClient(config)
	if err != nil {
		err = chkitErrors.ErrUnableToInitClient().
			AddDetailsErr(err)
		log.WithError(err).
			Error(err)
		return err
	}
	setClient(ctx, *client)
	return nil
}

func setupConfig(ctx *cli.Context) error {
	config := getConfig(ctx)
	err := loadConfig(ctx.String("config"), &config)
	if err != nil {
		return err
	}
	if ctx.Bool("test") {
		ctx.Set("api", testContainerumAPI)
	}
	if config.APIaddr == "" {
		config.APIaddr = ctx.String("api")
	}
	setConfig(ctx, config)
	return nil
}

func persist(ctx *cli.Context) error {
	if !ctx.IsSet("config") {
		return saveConfig(ctx)
	}
	return nil
}
