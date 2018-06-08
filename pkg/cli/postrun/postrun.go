package postrun

import (
	"fmt"

	"github.com/containerum/chkit/pkg/cli/login"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
)

func PostRun(ctx *context.Context) {
	var logger = ctx.Log.Component("PostRu")
	if ctx.Changed {
		err := configuration.SyncConfig(ctx)
		switch err {
		case nil:
			// pass
		case configuration.ErrIncompatibleConfig:
			ctx.Namespace = context.Namespace{}
			login.RunLogin(ctx, login.Flags{
				Username: ctx.Client.Username,
				Password: ctx.Client.Password,
			})
		default:
			logger.WithError(err).Errorf("unable to save config")
			fmt.Printf("Unable to save config: %v\n", err)
		}
	}
	if err := configuration.SaveTokens(ctx, ctx.Client.Tokens); err != nil {
		logger.WithError(err).Errorf("unable to save tokens")
		fmt.Printf("Unable to save tokens: %v\n", err)
	}
}
