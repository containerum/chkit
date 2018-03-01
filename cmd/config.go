package cmd

import (
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

func initConfig() error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	Configuration.ConfigPath = path.Join(home, configDir)
	return nil
}
