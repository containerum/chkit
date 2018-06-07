package clisetup

import (
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/sirupsen/logrus"
)

var Config = struct {
	DebugRequests bool
}{}

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
	logrus.Debugf("running setup")
	err := SetupConfig(ctx)
	switch {
	case err == nil:
		// pass
	case ErrInvalidUserInfo.Match(err):
		logrus.Debugf("invalid user information")
		logrus.Debugf("running login")
		if err := InteractiveLogin(ctx); err != nil {
			logrus.WithError(err).Errorf("unable to login")
			return err
		}
	default:
		logrus.WithError(ErrFatalError.Wrap(err)).Errorf("fatal error while config setup")
		return ErrFatalError.Wrap(err)
	}

	logrus.Debugf("client initialisation")
	if err := SetupClient(ctx, false); err != nil {
		logrus.WithError(err).Errorf("unable to init client")
		return err
	}
	if err := ctx.Client.Auth(); err != nil {
		logrus.WithError(err).Errorf("unable to auth")
		return err
	}

	logrus.Debugf("saving tokens")
	if err := configuration.SaveTokens(ctx, ctx.Client.Tokens); err != nil {
		logrus.WithError(err).Errorf("unable to save tokens")
		return err
	}

	if ctx.Namespace.IsEmpty() {
		return GetDefaultNS(ctx, false)
	}
	return nil
}
