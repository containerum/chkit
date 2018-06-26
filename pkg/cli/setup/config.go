package setup

import (
	"net/url"
	"os"

	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/util/fingerprint"
	"github.com/sirupsen/logrus"
)

func SetupConfig(ctx *context.Context) error {
	ctx.GetClient().Fingerprint = fingerpint.Fingerprint()
	tokens, err := configuration.LoadTokens(ctx)
	if err != nil && !os.IsNotExist(err) {
		return ErrUnableToLoadTokens.Wrap(err)
	} else if os.IsNotExist(err) {
		if err = configuration.SaveTokens(ctx, model.Tokens{}); err != nil {
			logrus.WithError(ErrUnableToSaveTokens.Wrap(err)).Errorf("unable to setup config")
			return ErrUnableToSaveTokens.Wrap(err)
		}
	}
	ctx.GetClient().Tokens = tokens
	if _, err := url.Parse(ctx.GetClient().APIaddr); err != nil {
		logrus.Debugf("invalid API url: %q", ctx.GetClient().APIaddr)
		return ErrInvalidAPIurl.Wrap(err)
	}
	if ctx.GetClient().Password == "" || ctx.GetClient().Username == "" {
		logrus.Debugf("invalid username or pass")
		return ErrInvalidUserInfo
	}
	return nil
}
