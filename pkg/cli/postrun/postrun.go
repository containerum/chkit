package postrun

import (
	"fmt"

	"os"

	"github.com/BurntSushi/toml"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/spf13/cobra"
)

func PostRunFunc(ctx *context.Context) func(*cobra.Command, []string) {
	return func(command *cobra.Command, strings []string) {
		PostRun(ctx)
	}
}

func PostRun(ctx *context.Context) {
	var logger = ctx.Log.Component("PostRun")
	logger.Debugf("START")
	defer logger.Debugf("END")
	if ctx.Changed {
		logger.Debugf("reading old config from %q", ctx.ConfigPath)
		oldConfig, err := configuration.ReadConfigFromDisk(ctx.ConfigPath)
		switch err {
		case nil, configuration.ErrIncompatibleConfig:
			logger.Struct(oldConfig)
			logger.Debugf("merging configs")
			ctx.SetStorable(ctx.GetStorable().Merge(oldConfig))
			logger.Struct(ctx.GetStorable())
		default:
			logger.WithError(err).Errorf("unable to read old config")
			angel.Angel(ctx, err)
		}
		logger.Debugf("writing new configuration")
		configFile, err := os.Create(ctx.ConfigPath)
		if err != nil {
			logger.Debugf("unable to create new config file")
			angel.Angel(ctx, err)
		}
		logger.Debugf("encoding configuration to TOML")
		if err := toml.NewEncoder(configFile).Encode(ctx.GetStorable()); err != nil {
			logger.Debugf("unable to encode configuration to TOML file")
			angel.Angel(ctx, err)
			os.Exit(1)
		}
	}
	logger.Debugf("saving tokens")
	if err := configuration.SaveTokens(ctx, ctx.Client.Tokens); err != nil {
		logger.WithError(err).Errorf("unable to save tokens")
		fmt.Printf("Unable to save tokens: %v\n", err)
	}
}
