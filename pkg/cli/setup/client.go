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

	//	ctx.SetSelfSignedTLSRule(certPolicy == AllowSelfSignedTLSCerts)

	if mode.DEBUG && !mode.MOCK {
		logger.Debugf("Using test API: %q", ctx.GetClient().APIaddr)
		//	ctx.GetClient().Log = logrus.StandardLogger().WriterLevel(logrus.DebugLevel)
		err = ctx.GetClient().Init(chClient.WithTestAPI)
	} else if mode.DEBUG && mode.MOCK {
		logger.Debugf("Using mock API")
		err = ctx.GetClient().Init(chClient.WithMock)
	} else if !mode.DEBUG {
		logger.Debugf("Using production API: %v", ctx.GetClient().APIaddr)
		err = ctx.GetClient().Init(chClient.WithCommonAPI)
	} else {
		panic(fmt.Sprintf("[setup.Client] invalid client mode state: DEBUG:%v MOCK:%v", mode.DEBUG, mode.MOCK))
	}
	if err != nil {
		logger.WithError(err).Errorf("unable to init client")
	}
	return nil
}
