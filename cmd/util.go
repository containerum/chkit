package cmd

import (
	"encoding/json"
	"os"
	"path"

	kubeClientModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/BurntSushi/toml"
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/model"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

func getLog(ctx *cli.Context) *logrus.Logger {
	return ctx.App.Metadata["log"].(*logrus.Logger)
}
func getConfig(ctx *cli.Context) model.ClientConfig {
	return ctx.App.Metadata["config"].(model.ClientConfig)
}

func setConfig(ctx *cli.Context, config model.ClientConfig) {
	ctx.App.Metadata["config"] = config
}

func saveConfig(ctx *cli.Context) {
	log := getLog(ctx)
	err := writeConfig(ctx)
	if err != nil {
		log.WithError(err).
			Fatalf("error while saving config")
	}
}

func getClient(ctx *cli.Context) chClient.Client {
	return ctx.App.Metadata["client"].(chClient.Client)
}

func setClient(ctx *cli.Context, client chClient.Client) {
	ctx.App.Metadata["client"] = client
}

func getConfigPath(ctx *cli.Context) string {
	return ctx.App.Metadata["configPath"].(string)
}
func exitOnErr(log *logrus.Logger, err error) {
	if err != nil {
		log.WithError(err).Fatal(err)
		os.Exit(1)
	}
}

func loadConfig(configFilePath string, config *model.ClientConfig) error {
	_, err := toml.DecodeFile(configFilePath, &config)
	if err != nil {
		return err
	}
	return nil
}

func writeConfig(ctx *cli.Context) error {
	configPath := getConfigPath(ctx)
	err := os.MkdirAll(configPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}
	file, err := os.Create(path.Join(configPath, "config.toml"))
	if err != nil {
		return err
	}
	config := getConfig(ctx)
	return toml.NewEncoder(file).Encode(config)
}

func saveTokens(ctx *cli.Context, tokens kubeClientModels.Tokens) error {
	file, err := os.Create(path.Join(getConfigPath(ctx), "tokens"))
	if err != nil {
		return err
	}
	return json.NewEncoder(file).Encode(tokens)
}

func loadTokens(config *model.Config) (kubeClientModels.Tokens, error) {
	tokens := kubeClientModels.Tokens{}
	file, err := os.Open(path.Join(config.ConfigPath, "tokens"))
	if err != nil {
		return tokens, err
	}
	return tokens, json.NewDecoder(file).Decode(&tokens)
}
