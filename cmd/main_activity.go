package cmd

import (
	"github.com/containerum/chkit/cmd/util"
	"gopkg.in/urfave/cli.v2"
)

func mainActivity(ctx *cli.Context) error {
	//client := getClient(ctx)
	log := util.GetLog(ctx)
	log.Infof("main activity")
	return nil
}
