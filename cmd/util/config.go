package util

import (
	"os"

	"github.com/BurntSushi/toml"
	. "github.com/containerum/chkit/cmd/context"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"
	cli "gopkg.in/urfave/cli.v2"
)

const (
	// ErrUnableToSaveConfig -- unable to save config
	ErrUnableToSaveConfig chkitErrors.Err = "unable to save config"
)

// GetConfigPath -- exctract config path from Context
func GetConfigPath(ctx *cli.Context) string {
	return ctx.App.Metadata["configPath"].(string)
}

// LoadConfig -- loads config from fs
func LoadConfig(configFilePath string, config *model.Config) error {
	_, err := toml.DecodeFile(configFilePath, &config.StorableConfig)
	if err != nil {
		return err
	}
	return nil
}

// SaveConfig -- writes config from Context to config dir
func SaveConfig() error {
	err := os.MkdirAll(Context.ConfigPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}
	file, err := os.Create(Context.ConfigPath)
	if err != nil {
		return err
	}
	return toml.NewEncoder(file).Encode(Context.ClientConfig.StorableConfig)
}
