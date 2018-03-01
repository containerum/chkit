package cmd

import (
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

func init() {
	home, err := homedir.Dir()
	if err != nil {
		log.WithError(err).
			Fatalf("error while getting homedir path")
	}
	Configuration.ConfigPath = path.Join(home, configDir)
}
