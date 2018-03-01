package cmd

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/containerum/chkit/pkg/model"
)

func exitOnErr(err error) {
	if err != nil {
		log.WithError(err).Fatal(err)
		os.Exit(1)
	}
}

func loadConfig(configFilePath string) error {
	config := model.Config{}
	_, err := toml.DecodeFile(configFilePath, &config)
	if err != nil {
		return err
	}
	Configuration = config
	return nil
}
