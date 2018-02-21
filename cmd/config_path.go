package cmd

import (
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

func init() {
	home, err := homedir.Dir()
	if err != nil {
		notepad.FATAL.Println(err)
	}
	Configuration.ConfigPath = path.Join(home, configDir)
}
