package cmd

import (
	"net/url"
	"os"

	"github.com/containerum/chkit/cmd/cmdutil"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/model"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

const (
	// ErrUnableToLoadConfig -- unable to load config
	ErrUnableToLoadConfig chkitErrors.Err = "unable to load config"
	// ErrInvalidUserInfo -- invalid user info"
	ErrInvalidUserInfo chkitErrors.Err = "invalid user info"
	// ErrInvalidAPIurl -- invalid API url
	ErrInvalidAPIurl chkitErrors.Err = "invalid API url"
	// ErrUnableToLoadTokens -- unable to load tokens
	ErrUnableToLoadTokens chkitErrors.Err = "unable to load tokens"
	// ErrUnableToSaveTokens -- unable to save tokens
	ErrUnableToSaveTokens chkitErrors.Err = "unable to save tokens"
)

func setupClient(ctx *cli.Context) error {
	config := cmdutil.GetConfig(ctx)
	var client *chClient.Client
	var err error
	config.APIaddr = API_ADDR

	if DEBUG && !MOCK {
		logrus.Debugf("Using test API: %q", config.APIaddr)
		if ctx.Bool("debug-requests") {
			logrus.Debugf("verbose requests logs")
			config.Log = logrus.StandardLogger().WriterLevel(logrus.DebugLevel)
		}
		client, err = chClient.NewClient(config, chClient.WithTestAPI)
	} else if DEBUG && MOCK {
		logrus.Debugf("Using mock API")
		client, err = chClient.NewClient(config, chClient.WithMock)
	} else {
		logrus.Debugf("Using production API: %v", config.APIaddr)
		client, err = chClient.NewClient(config, chClient.WithCommonAPI)
	}
	if err != nil {
		return err
	}
	cmdutil.SetClient(ctx, client)
	return nil
}

func setupConfig(ctx *cli.Context) error {
	config := cmdutil.GetConfig(ctx)
	logrus.Debugf("test: %q", ctx.String("test"))
	config.Fingerprint = Fingerprint()
	tokens, err := cmdutil.LoadTokens(ctx)
	if err != nil && !os.IsNotExist(err) {
		return ErrUnableToLoadTokens.Wrap(err)
	} else if os.IsNotExist(err) {
		if err = cmdutil.SaveTokens(ctx, model.Tokens{}); err != nil {
			return ErrUnableToSaveTokens.Wrap(err)
		}
	}
	config.Tokens = tokens
	if ctx.IsSet("test") {
		testAPIurl := os.Getenv("CONTAINERUM_API")
		logrus.Debugf("using test api %q", testAPIurl)
		config.APIaddr = testAPIurl
	}
	if _, err := url.Parse(config.APIaddr); err != nil {
		logrus.Debugf("invalid API url: %q", config.APIaddr)
		return ErrInvalidAPIurl.Wrap(err)
	}
	if config.Password == "" || config.Username == "" {
		logrus.Debugf("invalid username or pass")
		cmdutil.SetConfig(ctx, config)
		return ErrInvalidUserInfo
	}
	cmdutil.SetConfig(ctx, config)
	return nil
}

func persist(ctx *cli.Context) error {
	if !ctx.IsSet("config") {
		return cmdutil.SaveConfig(ctx)
	}
	return nil
}

func loadConfig(ctx *cli.Context) error {
	//log := cmdutil.GetLog(ctx)
	config := cmdutil.GetConfig(ctx)
	err := cmdutil.LoadConfig(ctx.String("config"), &config)
	if err != nil {
		return ErrUnableToLoadConfig.Wrap(err)
	}
	cmdutil.SetConfig(ctx, config)
	return nil
}

func setupAll(ctx *cli.Context) error {
	logrus.Debugf("setuping config")
	if err := loadConfig(ctx); err != nil {
		return err
	}
	if err := setupConfig(ctx); err != nil {
		return err
	}
	logrus.Debugf("setuping client")
	if err := setupClient(ctx); err != nil {
		return err
	}
	client := cmdutil.GetClient(ctx)
	logrus.Debugf("API: %q", client.APIaddr)
	return nil
}
