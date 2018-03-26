package cmd

import (
	"path"

	"github.com/containerum/chkit/cmd/config_dir"
	homedir "github.com/mitchellh/go-homedir"
)

func configPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return path.Join(home, confDir.ConfigDir), nil
}
