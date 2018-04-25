package clisetup

import (
	"net/url"
	"os"

	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/util/fingerprint"
	"github.com/sirupsen/logrus"
)

func SetupConfig() error {
	config := context.GlobalContext.Client.Config
	config.Fingerprint = fingerpint.Fingerprint()
	tokens, err := configuration.LoadTokens()
	if err != nil && !os.IsNotExist(err) {
		return ErrUnableToLoadTokens.Wrap(err)
	} else if os.IsNotExist(err) {
		if err = configuration.SaveTokens(model.Tokens{}); err != nil {
			logrus.WithError(ErrUnableToSaveTokens.Wrap(err)).Errorf("unable to setup config")
			return ErrUnableToSaveTokens.Wrap(err)
		}
	}
	config.Tokens = tokens
	if _, err := url.Parse(config.APIaddr); err != nil {
		logrus.Debugf("invalid API url: %q", config.APIaddr)
		return ErrInvalidAPIurl.Wrap(err)
	}
	if config.Password == "" || config.Username == "" {
		logrus.Debugf("invalid username or pass")
		return ErrInvalidUserInfo
	}
	context.GlobalContext.Client.Config = config
	return nil
}
