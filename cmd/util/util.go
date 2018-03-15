package util

import (
	"os"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

const (
	// ErrUnableToSaveConfig -- unable to save config
	ErrUnableToSaveConfig chkitErrors.Err = "unable to save config"
)

// GetLog -- extract logger instance from Context
func GetLog(ctx *cli.Context) *logrus.Logger {
	return ctx.App.Metadata["log"].(*logrus.Logger)
}

// ExitOnErr -- logs error and exit program
func ExitOnErr(log *logrus.Logger, err error) {
	if err != nil {
		log.WithError(err).Fatal(err)
		os.Exit(1)
	}
}