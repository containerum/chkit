package util

import (
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/containerum/chkit/pkg/model"
	cli "gopkg.in/urfave/cli.v2"
)

// GetConfig -- exctract config structure from Context
func GetConfig(ctx *cli.Context) model.Config {
	return ctx.App.Metadata["config"].(model.Config)
}

// SetConfig -- store config in Context
func SetConfig(ctx *cli.Context, config model.Config) {
	ctx.App.Metadata["config"] = config
}

// SaveConfig -- writes config in config path
func SaveConfig(ctx *cli.Context) error {
	err := WriteConfig(ctx)
	if err != nil {
		return ErrUnableToSaveConfig.Wrap(err)
	}
	return nil
}

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

// WriteConfig -- writes config from Context to config dir
func WriteConfig(ctx *cli.Context) error {
	configPath := GetConfigPath(ctx)
	err := os.MkdirAll(configPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}
	file, err := os.Create(path.Join(configPath, "config.toml"))
	if err != nil {
		return err
	}
	config := GetConfig(ctx)
	return toml.NewEncoder(file).Encode(config.StorableConfig)
}
