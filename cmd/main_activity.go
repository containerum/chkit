package cmd

import (
	"gopkg.in/urfave/cli.v2"
)

func mainActivity(ctx *cli.Context) error {
	client := getClient(ctx)
	if err := client.Login(); err != nil {
		return err
	}
	return nil
}
