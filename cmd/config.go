package cmd

import (
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

func getConfigPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return path.Join(home, configDir), nil
}
