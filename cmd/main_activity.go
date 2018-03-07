package cmd

import (
	"github.com/containerum/chkit/pkg/chkitErrors"
	"gopkg.in/urfave/cli.v2"
)

func mainActivity(ctx *cli.Context) error {
	client := getClient(ctx)
	if err := client.Login(); err != nil {
		return &chkitErrors.ExitCoder{
			Err:  err,
			Code: 2,
		}
	}
	setTokens(ctx, client.Tokens)
	if err := saveTokens(ctx, client.Tokens); err != nil {
		return &chkitErrors.ExitCoder{
			Err:  err,
			Code: 2,
		}
	}
	return nil
}
