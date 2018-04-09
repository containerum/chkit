package main

import (
	"os"

	"github.com/containerum/chkit/cmd"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

func main() {
	defer angel(recover())
	switch err := cmd.Run(os.Args).(type) {
	case nil:
		// pass
	case chkitErrors.Err, cli.ExitCoder:
		logrus.WithError(err).Error("fatal error")
	default:
		angel(err)
	}
}
