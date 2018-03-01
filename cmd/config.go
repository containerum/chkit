package cmd

import (
	"path"

	"github.com/containerum/chkit/pkg/model"
	homedir "github.com/mitchellh/go-homedir"
)

func initConfig(config *model.Config) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	config.ConfigPath = path.Join(home, configDir)
	return nil
}
