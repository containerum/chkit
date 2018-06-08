package configuration

import (
	"os"

	"github.com/containerum/chkit/pkg/context"

	"strings"

	"github.com/BurntSushi/toml"
	"github.com/containerum/chkit/pkg/chkitErrors"
)

const (
	// ErrUnableToSaveConfig -- unable to save config
	ErrUnableToSaveConfig chkitErrors.Err = "unable to sync config"
	// ErrUnableToLoadConfig -- unable to load config
	ErrUnableToLoadConfig chkitErrors.Err = "unable to load config"

	ErrIncompatibleConfig chkitErrors.Err = "it seems you try to run chkit with incompatible config file. " +
		"Please, delete config file and run 'chkit login'"
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
	config, err := ReadConfigFromDisk(ctx.ConfigPath)
	switch {
	case err == nil:
		// pass
	case ErrIncompatibleConfig == err:
		return err
	case !os.IsExist(err):
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

func ReadConfigFromDisk(configPath string) (context.Storable, error) {
	config := context.Storable{}
	_, err := toml.DecodeFile(configPath, &config)
	if err != nil && !os.IsNotExist(err) {
		if strings.Contains(err.Error(), " type mismatch") {
			return config, ErrIncompatibleConfig
		}
		return config, err
	}
	return config, nil
}
