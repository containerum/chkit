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
		logrus.WithField("operation", "SetupClient").Debugf("Using test API: %q", GlobalContext.Client.APIaddr)
		if Config.DebugRequests {
			logrus.Debugf("verbose requests logs")
			GlobalContext.Client.Log = logrus.StandardLogger().WriterLevel(logrus.DebugLevel)
		}
		err = GlobalContext.Client.Init(chClient.WithTestAPI)
	} else if mode.DEBUG && mode.MOCK {
		logrus.Debugf("Using mock API")
		err = GlobalContext.Client.Init(chClient.WithMock)
	} else {
		logrus.Debugf("Using production API: %v", GlobalContext.Client.APIaddr)
		err = GlobalContext.Client.Init(chClient.WithCommonAPI)
	}
	if err != nil {
		return err
	}
	return nil
}
