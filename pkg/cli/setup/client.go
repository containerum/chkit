package setup

import (
	"fmt"
	"net/url"
	"os"

	"github.com/containerum/chkit/pkg/cli/mode"
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/util/fingerprint"
	"github.com/sirupsen/logrus"
)

type CertPolicy string

const (
	DoNotAlloSelfSignedTLSCerts CertPolicy = ""
	AllowSelfSignedTLSCerts     CertPolicy = "allow self signed certs"
)

func (certPolicy CertPolicy) String() string {
	return string(certPolicy)
}

func Client(ctx *context.Context, certPolicy CertPolicy) error {
	var logger = ctx.Log.Component("client setup")
	logger.Debugf("START")
	defer logger.Debugf("END")

	ctx.Client.Fingerprint = fingerpint.Fingerprint()
	tokens, err := configuration.LoadTokens(ctx)
	if err != nil && !os.IsNotExist(err) {
		return ErrUnableToLoadTokens.Wrap(err)
	} else if os.IsNotExist(err) {
		if err = configuration.SaveTokens(ctx, model.Tokens{}); err != nil {
			logrus.WithError(ErrUnableToSaveTokens.Wrap(err)).Errorf("unable to setup config")
			return ErrUnableToSaveTokens.Wrap(err)
		}
	}
	ctx.Client.Tokens = tokens
	if _, err := url.Parse(ctx.Client.APIaddr); err != nil {
		logrus.Debugf("invalid API url: %q", ctx.Client.APIaddr)
		return ErrInvalidAPIurl.Wrap(err)
	}
	if ctx.Client.Password == "" || ctx.Client.Username == "" {
		logrus.Debugf("invalid username or pass")
		return ErrInvalidUserInfo
	}

	ctx.AllowSelfSignedTLS = certPolicy == AllowSelfSignedTLSCerts

	if mode.DEBUG && !mode.MOCK {
		logger.Debugf("Using test API: %q", ctx.Client.APIaddr)
		ctx.Client.Log = logrus.StandardLogger().WriterLevel(logrus.DebugLevel)
		err = ctx.Client.Init(chClient.WithTestAPI)
	} else if mode.DEBUG && mode.MOCK {
		logger.Debugf("Using mock API")
		err = ctx.Client.Init(chClient.WithMock)
	} else if !mode.DEBUG {
		logger.Debugf("Using production API: %v", ctx.Client.APIaddr)
		err = ctx.Client.Init(chClient.WithCommonAPI)
	} else {
		panic(fmt.Sprintf("[setup.Client] invalid client mode state: DEBUG:%v MOCK:%v", mode.DEBUG, mode.MOCK))
	}
	if err != nil {
		logger.WithError(err).Errorf("unable to init client")
	}
	return nil
}
