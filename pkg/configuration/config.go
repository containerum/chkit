package configuration

import (
	"os"

	"github.com/containerum/chkit/pkg/context"

	"github.com/BurntSushi/toml"
	"github.com/containerum/chkit/pkg/chkitErrors"
)

const (
	// ErrUnableToSaveConfig -- unable to save config
	ErrUnableToSaveConfig chkitErrors.Err = "unable to save config"
	// ErrUnableToLoadConfig -- unable to load config
	ErrUnableToLoadConfig chkitErrors.Err = "unable to load config"
)

// LoadConfig -- loads config from fs
func LoadConfig(ctx *context.Context) error {
	config := context.Storable{}
	_, err := toml.DecodeFile(ctx.ConfigPath, &config)
	if err != nil {
		return ErrUnableToLoadConfig.Wrap(err)
	}
	ctx.SetStorable(config)
	return nil
}

// SaveConfig -- writes config from Context to config dir
func SaveConfig(ctx *context.Context) error {
	err := os.MkdirAll(ctx.ConfigDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return ErrUnableToSaveConfig.Wrap(err)
	}
	file, err := os.Create(ctx.ConfigPath)
	if err != nil {
		return ErrUnableToSaveConfig.Wrap(err)
	}
	if err := toml.NewEncoder(file).Encode(ctx.GetStorable()); err != nil {
		return ErrUnableToSaveConfig.Wrap(err)
	}
	return nil
}
