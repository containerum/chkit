package postrun

import (
	"fmt"

	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/sirupsen/logrus"
)

func PostRun(logger logrus.FieldLogger, ctx *context.Context) {
	if ctx.Changed {
		if err := configuration.SyncConfig(ctx); err != nil {
			logger.WithError(err).Errorf("unable to save config")
			fmt.Printf("Unable to save config: %v\n", err)
		}
	}
	if err := configuration.SaveTokens(ctx, ctx.Client.Tokens); err != nil {
		logger.WithError(err).Errorf("unable to save tokens")
		fmt.Printf("Unable to save tokens: %v\n", err)
	}
}
