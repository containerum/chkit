package cmd

import (
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

var (
	ConfigPath string
)

func init() {
	home, err := homedir.Dir()
	if err != nil {
		notepad.FATAL.Println(err)
	}
	ConfigPath = path.Join(home, configDir)
}
