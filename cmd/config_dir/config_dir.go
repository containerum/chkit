package confDir

import (
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

func ConfigDir() string {
	home, err := homedir.Dir()
	if err != nil {
		panic("[config_dir ConfigDir] " + err.Error())
	}
	return path.Join(home, configDir)
}
