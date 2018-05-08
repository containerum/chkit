package clisetup

import (
	"github.com/containerum/chkit/pkg/cli/mode"
	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/context"
	"github.com/sirupsen/logrus"
)

func SetupClient(ctx *context.Context, debugRequests bool) error {
	var err error
	if (mode.DEBUG && !mode.MOCK) || ctx.AllowSelfSignedTLS {
		logrus.WithField("operation", "SetupClient").Debugf("Using test API: %q", ctx.Client.APIaddr)
		if debugRequests {
			logrus.Debugf("verbose requests logs")
			ctx.Client.Log = logrus.StandardLogger().WriterLevel(logrus.DebugLevel)
		}
		err = ctx.Client.Init(chClient.WithTestAPI)
	} else if mode.DEBUG && mode.MOCK {
		logrus.Debugf("Using mock API")
		err = ctx.Client.Init(chClient.WithMock)
	} else {
		logrus.Debugf("Using production API: %v", ctx.Client.APIaddr)
		err = ctx.Client.Init(chClient.WithCommonAPI)
	}
	if err != nil {
		return err
	}
	return nil
}
