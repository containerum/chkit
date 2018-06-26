package setup

import (
	"os"
	"path"
	"time"

	"github.com/containerum/chkit/pkg/configdir"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/sirupsen/logrus"
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
	ctx.Log = coblog.Log{logrus.StandardLogger()}
	return nil
}
