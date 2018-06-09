package coblog

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/oleiade/reflections"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	_ logrus.FieldLogger = Log{}

	Std = Log{logrus.StandardLogger()}
)

type Log struct {
	logrus.FieldLogger
}

func Component(component string, optionalLogger ...logrus.FieldLogger) Log {
	var logger logrus.FieldLogger
	if len(optionalLogger) > 0 {
		logger = optionalLogger[0]
	} else {
		logger = logrus.StandardLogger()
	}
	return Log{logger.WithField("component", component)}
}

func (log Log) Command(command string) Log {
	return Log{FieldLogger: log.FieldLogger.WithField("command", command)}
}

func (log Log) Component(component string) Log {
	return Log{FieldLogger: log.FieldLogger.WithField("component", component)}
}

func Logger(cmd *cobra.Command, optionalLogger ...logrus.FieldLogger) Log {
	var logger logrus.FieldLogger
	if len(optionalLogger) > 0 {
		logger = optionalLogger[0]
	} else {
		logger = logrus.StandardLogger()
	}
	return Log{logger.WithField(Field(cmd))}
}

func (log Log) StructFields(v interface{}) {
	var items, err = reflections.Items(v)
	if err != nil {
		log.WithError(err).Panicf("unable to encode value %v")
	}
	var structName = reflect.ValueOf(v).Type().Name()
	if structName == "" {
		structName = reflect.ValueOf(v).Kind().String()
	}
	var logger = log.WithField("data", structName)
	logger.Debugf("%v:", structName)
	var indent = strings.Repeat(" ", len(structName))
	for name, field := range items {
		logger.WithField(name, field).Infof(indent)
	}
}

func (log Log) Struct(v interface{}) {
	var items, err = reflections.Items(v)
	if err != nil {
		log.WithError(err).Panicf("unable to encode value %v")
	}
	var structName = reflect.ValueOf(v).Type().Name()
	var logger = log.WithField("data", structName)
	logger.Debugf("%v:", structName)
	var indent = strings.Repeat(" ", len(structName))
	for name, field := range items {
		logger.Debugf("%s%s : %v", indent, name, field)
	}
}

func Field(cmd *cobra.Command) (key string, value interface{}) {
	commandName := fmt.Sprintf("%s %s", cmd.Parent().Use, cmd.Use)
	return "command", commandName
}
