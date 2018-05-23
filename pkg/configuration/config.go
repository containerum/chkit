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

// SyncConfig -- writes config from Context to config dir
func SyncConfig(ctx *context.Context) error {
	err := os.MkdirAll(ctx.ConfigDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return ErrUnableToSaveConfig.Wrap(err)
	}
	config, err := loadConfig(ctx.ConfigPath)
	if err != nil && !os.IsExist(err) {
		return ErrUnableToSaveConfig.Wrap(err)
	}
	file, err := os.Create(ctx.ConfigPath)
	if err != nil {
		return ErrUnableToSaveConfig.Wrap(err)
	}
	if ctx.Changed {
		config = config.Merge(ctx.GetStorable())
	} else {
		config = ctx.GetStorable().Merge(config)
	}
	ctx.SetStorable(config)
	if err := toml.NewEncoder(file).Encode(config); err != nil {
		return ErrUnableToSaveConfig.Wrap(err)
	}
	return nil
}

func loadConfig(configPath string) (context.Storable, error) {
	config := context.Storable{}
	_, err := toml.DecodeFile(configPath, &config)
	if err != nil && !os.IsNotExist(err) {
		return config, err
	}
	return config, nil
}
