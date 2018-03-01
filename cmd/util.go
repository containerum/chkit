package cmd

import (
	"encoding/json"
	"os"
	"path"

	kubeClientModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/BurntSushi/toml"
	"github.com/containerum/chkit/pkg/model"
	"github.com/sirupsen/logrus"
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

func saveTokens(config *model.Config) error {
	file, err := os.Create(path.Join(config.ConfigPath, "tokens"))
	if err != nil {
		return err
	}
	return json.NewEncoder(file).Encode(&config.Tokens)
}

func loadTokens(config *model.Config) (kubeClientModels.Tokens, error) {
	tokens := kubeClientModels.Tokens{}
	file, err := os.Open(path.Join(config.ConfigPath, "tokens"))
	if err != nil {
		return tokens, err
	}
	return tokens, json.NewDecoder(file).Decode(&tokens)
}
