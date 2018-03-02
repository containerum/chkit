package cmd

import (
	"os"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/client"

	"github.com/containerum/chkit/pkg/model"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

func setupClient(log *logrus.Logger, ctx *cli.Context) error {
	clientConfig := model.ClientConfig{}
	err := loadConfig(ctx.String("config"), &clientConfig)
	if clientConfig.APIaddr == "" {
		clientConfig.APIaddr = ctx.String("api")
	}
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
	}
	err = saveConfig(getConfigPath(ctx), &clientConfig)
	if err != nil {
		log.WithError(err).
			Errorf("error while saving config")
		return err
	}
	client, err := chClient.NewClient(clientConfig)
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
