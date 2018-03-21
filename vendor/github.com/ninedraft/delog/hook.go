package delog

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var (
	_ logrus.Hook = new(Hook)
)

// Hook -- logrus hook. Writes log messages to provided writer with prefix [FILE_PATH:FUNCTION:LINE]
type Hook struct {
	writer    io.Writer
	formatter Formatter
	logLevels []logrus.Level
}

// NewHook -- creates Hook with provided logrus Formatter, log writer and optional log Levels (debug by default). Adds [FILE_PATH:FUNCTION:LINE] prefix on Fire
func NewHook(formatter logrus.Formatter, wr io.Writer, logLevels ...logrus.Level) *Hook {
	if len(logLevels) == 0 {
		logLevels = []logrus.Level{logrus.DebugLevel}
	}
	if formatter == nil {
		formatter = &logrus.TextFormatter{FullTimestamp: true}
	}
	if wr == nil {
		wr = os.Stdout
	}
	delogFormatter := *NewFormatter(formatter)
	delogFormatter.stackOffset = 2
	return &Hook{
		writer:    wr,
		formatter: delogFormatter,
		logLevels: logLevels,
	}
}

// Levels -- returns log levels, consumed by logrus
func (hook *Hook) Levels() []logrus.Level {
	return hook.logLevels
}

// Fire -- hook method, is called by logrus
func (hook *Hook) Fire(entry *logrus.Entry) error {
	msg, err := hook.formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.writer.Write(msg)
	return err
}
