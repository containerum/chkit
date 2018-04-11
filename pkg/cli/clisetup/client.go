package clisetup

import (
	"github.com/containerum/chkit/pkg/cli/mode"
	"github.com/containerum/chkit/pkg/client"
	. "github.com/containerum/chkit/pkg/context"
	"github.com/sirupsen/logrus"
)

func SetupClient() error {
	var err error
	if mode.DEBUG && !mode.MOCK {
		logrus.WithField("operation", "SetupClient").Debugf("Using test API: %q", Context.Client.APIaddr)
		if Config.DebugRequests {
			logrus.Debugf("verbose requests logs")
			Context.Client.Log = logrus.StandardLogger().WriterLevel(logrus.DebugLevel)
		}
		err = Context.Client.Init(chClient.WithTestAPI)
	} else if mode.DEBUG && mode.MOCK {
		logrus.Debugf("Using mock API")
		err = Context.Client.Init(chClient.WithMock)
	} else {
		logrus.Debugf("Using production API: %v", Context.Client.APIaddr)
		err = Context.Client.Init(chClient.WithCommonAPI)
	}
	if err != nil {
		return err
	}
	return nil
}
