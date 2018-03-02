package cmd

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/model"

	"gopkg.in/urfave/cli.v2"
)

func setupClient(ctx *cli.Context) error {
	log := getLog(ctx)
	client, err := chClient.NewClient(getConfig(ctx))
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

func setup(ctx *cli.Context) error {
	log := getLog(ctx)
	config := model.ClientConfig{}
	err := loadConfig(ctx.String("config"), &config)
	if err != nil && !os.IsNotExist(err) {
		log.WithError(err).
			Errorf("error while loading config file")
		return err
	} else if os.IsNotExist(err) {
		log.Info("You are not logged in!")
		err = ctx.App.Command("login").Run(ctx)
		if err != nil {
			return err
		}
		config = getConfig(ctx)
		fmt.Println(config)
	}
	if config.APIaddr == "" {
		config.APIaddr = ctx.String("api")
	}
	setConfig(ctx, config)
	if !ctx.IsSet("config") {
		saveConfig(ctx, config)
	}
	return nil
}
