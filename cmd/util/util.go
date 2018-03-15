package util

import (
	"encoding/json"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/model"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

const (
	ErrUnableToSaveConfig chkitErrors.Err = "unable to save config"
)

func GetLog(ctx *cli.Context) *logrus.Logger {
	return ctx.App.Metadata["log"].(*logrus.Logger)
}
func GetConfig(ctx *cli.Context) model.Config {
	return ctx.App.Metadata["config"].(model.Config)
}

func SetConfig(ctx *cli.Context, config model.Config) {
	ctx.App.Metadata["config"] = config
}

func SaveConfig(ctx *cli.Context) error {
	err := WriteConfig(ctx)
	if err != nil {
		return ErrUnableToSaveConfig.Wrap(err)
	}
	return nil
}

func GetClient(ctx *cli.Context) chClient.Client {
	return ctx.App.Metadata["client"].(chClient.Client)
}

func SetClient(ctx *cli.Context, client chClient.Client) {
	ctx.App.Metadata["client"] = client
}

func GetConfigPath(ctx *cli.Context) string {
	return ctx.App.Metadata["configPath"].(string)
}
func ExitOnErr(log *logrus.Logger, err error) {
	if err != nil {
		log.WithError(err).Fatal(err)
		os.Exit(1)
	}
}

func LoadConfig(configFilePath string, config *model.Config) error {
	_, err := toml.DecodeFile(configFilePath, &config.UserInfo)
	if err != nil {
		return err
	}
	return nil
}

func WriteConfig(ctx *cli.Context) error {
	configPath := GetConfigPath(ctx)
	err := os.MkdirAll(configPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}
	file, err := os.Create(path.Join(configPath, "config.toml"))
	if err != nil {
		return err
	}
	config := GetConfig(ctx)
	return toml.NewEncoder(file).Encode(config.UserInfo)
}

func GetTokens(ctx *cli.Context) model.Tokens {
	return ctx.App.Metadata["tokens"].(model.Tokens)
}

func SetTokens(ctx *cli.Context, tokens model.Tokens) {
	ctx.App.Metadata["tokens"] = tokens
}
func SaveTokens(ctx *cli.Context, tokens model.Tokens) error {
	file, err := os.Create(path.Join(GetConfigPath(ctx), "tokens"))
	if err != nil {
		return err
	}
	return json.NewEncoder(file).Encode(tokens)
}

func LoadTokens(ctx *cli.Context) (model.Tokens, error) {
	tokens := model.Tokens{}
	file, err := os.Open(path.Join(GetConfigPath(ctx), "tokens"))
	if err != nil {
		return tokens, err
	}
	return tokens, json.NewDecoder(file).Decode(&tokens)
}
