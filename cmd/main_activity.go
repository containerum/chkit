package cmd

import (
	"gopkg.in/urfave/cli.v2"
)

func mainActivity(ctx *cli.Context) error {
	//client := getClient(ctx)
	log := getLog(ctx)
	log.Infof("main activity")
	return nil
}
