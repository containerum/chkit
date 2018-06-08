package login

import (
	"github.com/containerum/chkit/pkg/cli/clisetup"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/util/coblog"
)

func Setup(ctx *context.Context) error {
	var logger = coblog.Component("login setup client")
	logger.Debugf("running login client setup")
	defer logger.Debugf("end login client setup")
	err := clisetup.SetupConfig(ctx)
	switch {
	case err == nil:
		// pass
	case clisetup.ErrInvalidUserInfo.Match(err):
		logger.Debugf("invalid user information")
		logger.Debugf("running login")
		if err := clisetup.InteractiveLogin(ctx); err != nil {
			logger.WithError(err).Errorf("unable to login")
			return err
		}
	default:
		logger.WithError(clisetup.ErrFatalError.Wrap(err)).Errorf("fatal error while config Setup")
		return clisetup.ErrFatalError.Wrap(err)
	}
	ctx.Client.Tokens = model.Tokens{}
	logger.Debugf("client initialisation")
	if err := clisetup.SetupClient(ctx, false); err != nil {
		logger.WithError(err).Errorf("unable to init client")
		return err
	}
	if err := ctx.Client.Auth(); err != nil {
		logger.WithError(err).Errorf("unable to auth")
		return err
	}

	logger.Debugf("saving tokens")
	if err := configuration.SaveTokens(ctx, ctx.Client.Tokens); err != nil {
		logger.WithError(err).Errorf("unable to save tokens")
		return err
	}
	return nil
}
