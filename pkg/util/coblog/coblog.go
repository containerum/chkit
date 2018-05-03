package coblog

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Logger(cmd *cobra.Command, optionalLogger ...logrus.FieldLogger) logrus.FieldLogger {
	var logger logrus.FieldLogger
	if len(optionalLogger) > 0 {
		logger = optionalLogger[0]
	} else {
		logger = logrus.StandardLogger()
	}
	return logger.WithField(Field(cmd))
}

func Field(cmd *cobra.Command) (key string, value interface{}) {
	commandName := fmt.Sprintf("%s %s", cmd.Parent().Use, cmd.Use)
	return "command", commandName
}
