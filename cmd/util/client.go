package util

import (
	"github.com/containerum/chkit/pkg/client"
	cli "gopkg.in/urfave/cli.v2"
)

func GetClient(ctx *cli.Context) chClient.Client {
	return ctx.App.Metadata["client"].(chClient.Client)
}

func SetClient(ctx *cli.Context, client chClient.Client) {
	ctx.App.Metadata["client"] = client
}
