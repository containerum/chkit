package util

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	_ logrus.Formatter = new(LogDebugger)
)

type LogDebugger struct {
	formatter  logrus.Formatter
	stackLevel int
	logLevels  []logrus.Level
}

func NewLogDebugger(stackLevel uint, formatter logrus.Formatter, logLevels ...logrus.Level) *LogDebugger {
	if len(logLevels) == 0 {
		logLevels = []logrus.Level{logrus.DebugLevel}
	}
	if formatter == nil {
		formatter = &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006 Jan Mon 15:04:05",
		}
	}
	return &LogDebugger{
		formatter:  formatter,
		logLevels:  logLevels,
		stackLevel: int(stackLevel),
	}
}

func (logdbg *LogDebugger) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Message = fmt.Sprintf("[%s] %s",
		debugData(logdbg.stackLevel+1),
		entry.Message)
	msg, err := logdbg.formatter.Format(entry)
	return msg, err
}

func debugData(stackLevel int) string {
	caller, filePath, line, _ := runtime.Caller(stackLevel)
	frame, _ := runtime.CallersFrames([]uintptr{caller}).Next()
	file := path.Base(path.Dir(filePath)) + "/" + path.Base(filePath)

	fnName := strings.Split(path.Base(frame.Function), ".")[1]
	return fmt.Sprintf("%d++%s:%s:%d", stackLevel, file, fnName, line)
}

func DebugData() string {
	return debugData(2)
}
