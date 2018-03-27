package confDir

import (
	"os"
	"path"
	"path/filepath"

	"github.com/containerum/chkit/pkg/chkitErrors"
	homedir "github.com/mitchellh/go-homedir"
)

const (
	ErrUnableToCreateConfigDir  chkitErrors.Err = "unable to create config dir"
	ErrUnableToCreateConfigFile chkitErrors.Err = "unable to create config file"
)

func ConfigDir() string {
	home, err := homedir.Dir()
	if err != nil {
		panic("[config_dir ConfigDir] " + err.Error())
	}
	home = filepath.ToSlash(home)
	return path.Join(home, configDir)
}

func init() {
	pathToConfigDir := ConfigDir()
	err := os.MkdirAll(pathToConfigDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		panic(ErrUnableToCreateConfigDir.Wrap(err))
	}

	pathToConfigFile := path.Join(pathToConfigDir, "config.toml")
	_, err = os.Stat(pathToConfigFile)
	if err != nil && os.IsNotExist(err) {
		file, err := os.Create(pathToConfigFile)
		if err != nil {
			panic(ErrUnableToCreateConfigFile.Wrap(err))
		}
		if err = file.Close(); err != nil {
			panic(ErrUnableToCreateConfigFile.Wrap(err))
		}
	} else if err != nil {
		panic(ErrUnableToCreateConfigFile.Wrap(err))
	}
}
