package cmdutil

import (
	"github.com/containerum/chkit/pkg/client"
	cli "gopkg.in/urfave/cli.v2"
)

// GetClient -- extract chkit Client
func GetClient(ctx *cli.Context) *chClient.Client {
	client := ctx.App.Metadata["client"].(chClient.Client)
	return &client
}

// SetClient -- store Client in Context
func SetClient(ctx *cli.Context, client *chClient.Client) {
	ctx.App.Metadata["client"] = *client
}
