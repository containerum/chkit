package prerun

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/cli/clisetup"
	"github.com/containerum/chkit/pkg/cli/login"
	"github.com/containerum/chkit/pkg/configdir"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/sirupsen/logrus"
)

const (
	// ErrFatalError -- unrecoverable fatal error
	ErrFatalError chkitErrors.Err = "fatal error"
)

func SetupLogs(ctx *context.Context) error {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC1123,
	})
	logFile := path.Join(configdir.LogDir(), configuration.LogFileName())
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logrus.Fatalf("error while creating log file: %v", err)
	}
	logrus.SetOutput(file)
	return nil
}

func PreRun(ctx *context.Context) error {
	SetupLogs(ctx)
	logrus.Debugf("loading config")
	if err := configuration.LoadConfig(ctx); err != nil {
		logrus.WithError(err).Errorf("unable to load config")
		return err
	}

	logrus.Debugf("running setup")
	err := clisetup.SetupConfig(ctx)
	switch {
	case err == nil:
		// pass
	case clisetup.ErrInvalidUserInfo.Match(err):
		logrus.Debugf("invalid user information")
		logrus.Debugf("running login")
		if err := login.InteractiveLogin(ctx); err != nil {
			logrus.WithError(err).Errorf("unable to login")
			return err
		}
	default:
		logrus.WithError(ErrFatalError.Wrap(err)).Errorf("fatal error while config setup")
		return ErrFatalError.Wrap(err)
	}

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

	if ctx.Namespace == "" {
		logrus.Debugf("getting user namespaces list")
		list, err := ctx.Client.GetNamespaceList()
		if err != nil {
			logrus.WithError(err).Errorf("unable to get user namespace list")
			fmt.Printf("Unable to get default namespace\n")
			return err
		}
		if len(list) == 0 {
			fmt.Printf("You have no namespaces!\n")
		} else {
			ctx.Changed = true
			ctx.Namespace = list[0].Label
		}
	}
	return nil
}
