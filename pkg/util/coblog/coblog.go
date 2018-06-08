package coblog

import (
	"fmt"

	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	_ logrus.FieldLogger = Log{}
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

func (log Log) Struct(v interface{}) {
	var value = reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return
	}
	var tt = value.Type()
	for fieldIndex := 0; fieldIndex < value.NumField(); fieldIndex++ {
		var name = tt.Field(fieldIndex).Name
		var field = value.Field(fieldIndex)
		if !field.CanSet() {
			continue
		}
		log.WithField("struct", tt.Name()).Printf("%s : %v", name, field.Interface())
	}
}

func Field(cmd *cobra.Command) (key string, value interface{}) {
	commandName := fmt.Sprintf("%s %s", cmd.Parent().Use, cmd.Use)
	return "command", commandName
}
