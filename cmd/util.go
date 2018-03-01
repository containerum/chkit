package cmd

import (
	"os"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/BurntSushi/toml"
	"github.com/containerum/chkit/pkg/model"
)

func exitOnErr(log *logrus.Logger, err error) {
	if err != nil {
		log.WithError(err).Fatal(err)
		os.Exit(1)
	}
}

func loadConfig(config *model.ClientConfig, configFilePath string) error {
	_, err := toml.DecodeFile(configFilePath, &config)
	if err != nil {
		return err
	}
	return nil
}

func saveConfig(config *model.Config) error {
	err := os.MkdirAll(config.ConfigPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}
	file, err := os.Create(path.Join(config.ConfigPath, "config.toml"))
	if err != nil {
		return err
	}
	return toml.NewEncoder(file).Encode(&config.Client)
}
