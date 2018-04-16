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
func LoadConfig() error {
	config := context.Storable{}
	_, err := toml.DecodeFile(context.GlobalContext.ConfigPath, &config)
	if err != nil {
		return ErrUnableToLoadConfig.Wrap(err)
	}
	context.GlobalContext.SetStorable(config)
	return nil
}

// SaveConfig -- writes config from Context to config dir
func SaveConfig() error {
	err := os.MkdirAll(context.GlobalContext.ConfigDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return ErrUnableToSaveConfig.Wrap(err)
	}
	file, err := os.Create(context.GlobalContext.ConfigPath)
	if err != nil {
		return ErrUnableToSaveConfig.Wrap(err)
	}
	if err := toml.NewEncoder(file).Encode(context.GlobalContext.GetStorable()); err != nil {
		return ErrUnableToSaveConfig.Wrap(err)
	}
	return nil
}
