package cmd

import (
	"encoding/json"
	"os"
	"path"

	"github.com/containerum/chkit/pkg/chkitErrors"

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
func getConfig(ctx *cli.Context) model.Config {
	return ctx.App.Metadata["config"].(model.Config)
}

func setConfig(ctx *cli.Context, config model.Config) {
	ctx.App.Metadata["config"] = config
}

func saveConfig(ctx *cli.Context) error {
	err := writeConfig(ctx)
	if err != nil {
		return chkitErrors.ErrUnableToSaveConfig().
			AddDetailsErr(err)
	}
	return nil
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

func loadConfig(configFilePath string, config *model.Config) error {
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

func getTokens(ctx *cli.Context) kubeClientModels.Tokens {
	return ctx.App.Metadata["tokens"].(kubeClientModels.Tokens)
}

func setTokens(ctx *cli.Context, tokens kubeClientModels.Tokens) {
	ctx.App.Metadata["tokens"] = tokens
}
func saveTokens(ctx *cli.Context, tokens kubeClientModels.Tokens) error {
	file, err := os.Create(path.Join(getConfigPath(ctx), "tokens"))
	if err != nil {
		return err
	}
	return json.NewEncoder(file).Encode(tokens)
}

func loadTokens(ctx *cli.Context) (kubeClientModels.Tokens, error) {
	tokens := kubeClientModels.Tokens{}
	file, err := os.Open(path.Join(getConfigPath(ctx), "tokens"))
	if err != nil {
		return tokens, err
	}
	return tokens, json.NewDecoder(file).Decode(&tokens)
}
