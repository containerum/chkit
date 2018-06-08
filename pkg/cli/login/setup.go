package login

import (
	"github.com/containerum/chkit/pkg/cli/clisetup"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/sirupsen/logrus"
)

func Setup(ctx *context.Context) error {
	logrus.Debugf("running Setup")
	err := clisetup.SetupConfig(ctx)
	switch {
	case err == nil:
		// pass
	case clisetup.ErrInvalidUserInfo.Match(err):
		logrus.Debugf("invalid user information")
		logrus.Debugf("running login")
		if err := clisetup.InteractiveLogin(ctx); err != nil {
			logrus.WithError(err).Errorf("unable to login")
			return err
		}
	default:
		logrus.WithError(clisetup.ErrFatalError.Wrap(err)).Errorf("fatal error while config Setup")
		return clisetup.ErrFatalError.Wrap(err)
	}
	ctx.Client.Tokens = model.Tokens{}
	logrus.Debugf("client initialisation")
	if err := clisetup.SetupClient(ctx, false); err != nil {
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
	return nil
}
