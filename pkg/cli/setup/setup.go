package setup

import (
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/util/coblog"
)

const (
	// ErrFatalError -- unrecoverable fatal error
	ErrFatalError chkitErrors.Err = "fatal error"
	// ErrInvalidUserInfo -- invalid user info"
	ErrInvalidUserInfo chkitErrors.Err = "invalid user info"
	// ErrInvalidAPIurl -- invalid API url
	ErrInvalidAPIurl chkitErrors.Err = "invalid API url"
	// ErrUnableToLoadTokens -- unable to load tokens
	ErrUnableToLoadTokens chkitErrors.Err = "unable to load tokens"
	// ErrUnableToSaveTokens -- unable to save tokens
	ErrUnableToSaveTokens chkitErrors.Err = "unable to save tokens"
)

func Setup(ctx *context.Context) error {
	var logger = coblog.Component("login setup client")
	logger.Debugf("running login client setup")
	defer logger.Debugf("end login client setup")
	err := SetupConfig(ctx)
	switch {
	case err == nil:
		// pass
	case ErrInvalidUserInfo.Match(err):
		logger.Debugf("invalid user information")
		logger.Debugf("running login")
		if err := InteractiveLogin(ctx); err != nil {
			logger.WithError(err).Errorf("unable to login")
			return err
		}
	default:
		logger.WithError(ErrFatalError.Wrap(err)).Errorf("fatal error while config Setup")
		return ErrFatalError.Wrap(err)
	}
	ctx.Client.Tokens = model.Tokens{}
	logger.Debugf("client initialisation")
	if err := SetupClient(ctx, false); err != nil {
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
