package configdir

import (
	"os"
	"path"
	"path/filepath"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/mitchellh/go-homedir"
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

func LogDir() string {
	home, err := homedir.Dir()
	if err != nil {
		panic("[config_dir LogDir] " + err.Error())
	}
	home = filepath.ToSlash(home)
	return path.Join(home, logDir)
}

func init() {
	pathToConfigDir := ConfigDir()
	err := os.MkdirAll(pathToConfigDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		panic(ErrUnableToCreateConfigDir.Wrap(err))
	}

	err = os.MkdirAll(LogDir(), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		panic(ErrUnableToCreateConfigDir.Wrap(err))
	}

	pathToConfigFile := path.Join(pathToConfigDir, "config.toml")
	file, err := os.OpenFile(pathToConfigFile, os.O_CREATE, 0600)
	switch {
	case err == nil || os.IsExist(err):
		os.Chmod(file.Name(), 0600)
		file.Close()
	default:
		panic(ErrUnableToCreateConfigFile.Wrap(err))
	}
}
