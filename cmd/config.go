package cmd

import (
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

func configPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return path.Join(home, configDir), nil
}
