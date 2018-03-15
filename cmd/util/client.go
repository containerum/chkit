package util

import (
	"github.com/containerum/chkit/pkg/client"
	cli "gopkg.in/urfave/cli.v2"
)

// GetClient -- extract chkit Client
func GetClient(ctx *cli.Context) chClient.Client {
	return ctx.App.Metadata["client"].(chClient.Client)
}

// SetClient -- store Client in Context
func SetClient(ctx *cli.Context, client chClient.Client) {
	ctx.App.Metadata["client"] = client
}
