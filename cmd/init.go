package cmd

import (
	"github.com/sirupsen/logrus"

	chClient "github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/model"
	cli "gopkg.in/urfave/cli.v2"
)

const (
	FlagAPIaddr    = "apiaddr"
	FlagConfigFile = "config"
)

var Configuration = struct {
	ConfigPath   string
	ConfigFile   string
	TokenFile    string
	ClientConfig model.Config
}{}

var ChkitClient chClient.Client
var log = &logrus.Logger{
	Formatter: &logrus.TextFormatter{},
}

var App = &cli.App{
	Name: "chkit",
	Action: func(ctx *cli.Context) error {
		return nil
	},
}
