package cmd

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

func mainActivity(ctx *cli.Context) error {
	//client := getClient(ctx)
	logrus.Infof("main activity")
	return nil
}
