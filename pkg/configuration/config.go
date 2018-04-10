package configuration

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/containerum/chkit/pkg/chkitErrors"
	. "github.com/containerum/chkit/pkg/context"
)

const (
	// ErrUnableToSaveConfig -- unable to save config
	ErrUnableToSaveConfig chkitErrors.Err = "unable to save config"
	// ErrUnableToLoadConfig -- unable to load config
	ErrUnableToLoadConfig chkitErrors.Err = "unable to load config"
)

// LoadConfig -- loads config from fs
func LoadConfig() error {
	_, err := toml.DecodeFile(Context.ConfigPath, &Context.Client.StorableConfig)
	if err != nil {
		return ErrUnableToLoadConfig.Wrap(err)
	}
	return nil
}

// SaveConfig -- writes config from Context to config dir
func SaveConfig() error {
	err := os.MkdirAll(Context.ConfigPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return ErrUnableToSaveConfig.Wrap(err)
	}
	file, err := os.Create(Context.ConfigPath)
	if err != nil {
		return ErrUnableToSaveConfig.Wrap(err)
	}
	if err := toml.NewEncoder(file).Encode(Context.Client.StorableConfig); err != nil {
		return ErrUnableToSaveConfig.Wrap(err)
	}
	return nil
}
