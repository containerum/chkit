package cmd

import (
	"github.com/containerum/chkit/pkg/chkitErrors"
	"gopkg.in/urfave/cli.v2"
)

func mainActivity(ctx *cli.Context) error {
	client := getClient(ctx)
	log := getLog(ctx)
	if err := client.Login(); err != nil {
		return &chkitErrors.ExitCoder{
			Err:  err,
			Code: 2,
		}
	}
	log.Infof("Trying to auth...")
	setTokens(ctx, client.Tokens)
	if err := saveTokens(ctx, client.Tokens); err != nil {
		return &chkitErrors.ExitCoder{
			Err:  err,
			Code: 2,
		}
	}
	log.Infof("OK")
	return nil
}
