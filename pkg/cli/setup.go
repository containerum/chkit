package cli

import (
	"net/url"
	"os"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/configuration"
	. "github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/util/fingerprint"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

const (

	// ErrInvalidUserInfo -- invalid user info"
	ErrInvalidUserInfo chkitErrors.Err = "invalid user info"
	// ErrInvalidAPIurl -- invalid API url
	ErrInvalidAPIurl chkitErrors.Err = "invalid API url"
	// ErrUnableToLoadTokens -- unable to load tokens
	ErrUnableToLoadTokens chkitErrors.Err = "unable to load tokens"
	// ErrUnableToSaveTokens -- unable to save tokens
	ErrUnableToSaveTokens chkitErrors.Err = "unable to save tokens"
)

func setupClient() error {
	var client *chClient.Client
	var err error
	if DEBUG && !MOCK {
		logrus.Debugf("Using test API: %q", Context.APIaddr)
		if runContext.DebugRequests {
			logrus.Debugf("verbose requests logs")
			Context.ClientConfig.Log = logrus.StandardLogger().WriterLevel(logrus.DebugLevel)
		}
		client, err = chClient.NewClient(Context.ClientConfig, chClient.WithTestAPI)
	} else if DEBUG && MOCK {
		logrus.Debugf("Using mock API")
		client, err = chClient.NewClient(Context.ClientConfig, chClient.WithMock)
	} else {
		logrus.Debugf("Using production API: %v", Context.APIaddr)
		client, err = chClient.NewClient(Context.ClientConfig, chClient.WithCommonAPI)
	}
	if err != nil {
		return err
	}
	Context.Client = client
	return nil
}

func setupConfig() error {
	config := Context.ClientConfig
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
	Context.ClientConfig = config
	Context.ClientConfig.APIaddr = Context.APIaddr
	return nil
}

func setupAll(ctx *cli.Context) error {
	logrus.Debugf("setuping config")
	if err := configuration.LoadConfig(); err != nil {
		return err
	}
	if err := setupConfig(); err != nil {
		return err
	}
	logrus.Debugf("setuping client")
	if err := setupClient(); err != nil {
		return err
	}
	logrus.Debugf("API: %q", Context.Client.APIaddr)
	return nil
}
