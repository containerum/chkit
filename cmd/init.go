package cmd

import (
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/model"
	jww "github.com/spf13/jwalterweatherman"
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

var ChkitClient client.ChkitClient
var notepad *jww.Notepad

var App = &cli.App{
	Name: "chkit",
	Action: func(ctx *cli.Context) error {
		return nil
	},
}
