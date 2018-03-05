package cmd

import (
	"fmt"

	"gopkg.in/urfave/cli.v2"
)

func mainActivity(ctx *cli.Context) error {
	fmt.Println("main activity")
	client := getClient(ctx)
	if err := client.Login(); err != nil {
		return err
	}
	setTokens(ctx, client.Tokens)
	saveTokens(ctx, client.Tokens)
	return nil
}
