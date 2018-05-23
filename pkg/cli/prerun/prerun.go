package prerun

import (
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/cli/clisetup"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/sirupsen/logrus"
)

const (
	// ErrFatalError -- unrecoverable fatal error
	ErrFatalError chkitErrors.Err = "fatal error"
)

func PreRun(ctx *context.Context) error {
	clisetup.SetupLogs(ctx)
	logrus.Debugf("loading config")
	if err := configuration.SyncConfig(ctx); err != nil {
		logrus.WithError(err).Errorf("unable to load config")
		return err
	}
	return clisetup.Setup(ctx)
}
