package delog

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strings"
)

var (
	_ logrus.Formatter = new(Formatter)
)

// Formatter -- logrus formatter, adds prefix a-la [FILE_PATH:FUNCTION:LINE]
type Formatter struct {
	stackOffset int
	formatter   logrus.Formatter
}

// NewFormatter -- creates new delog logrus formatter with provided underlying formatter. If formatter param is nil, them uses default text formatter
func NewFormatter(formatter logrus.Formatter) *Formatter {
	if formatter == nil {
		formatter = &logrus.TextFormatter{}
	}
	return &Formatter{
		stackOffset: 0,
		formatter:   formatter,
	}
}

// Format -- add prefix [FILE_PATH:FUNCTION:LINE] to entry message and pass it to underlying formatter
func (dbgFormatter *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	skip := 0
	if len(entry.Data) == 0 {
		skip = 6
	} else {
		skip = 4
	}
	skip += dbgFormatter.stackOffset

	caller, filePath, line, _ := runtime.Caller(skip)
	frame, _ := runtime.CallersFrames([]uintptr{caller}).Next()
	file := stripGopath(filePath)
	fnName := strings.Split(path.Base(frame.Function), ".")[1]

	ddbgData := fmt.Sprintf("%s:%s:%d", file, fnName, line)
	entry.Message = fmt.Sprintf("[%s] %s", ddbgData, entry.Message)
	msg, err := dbgFormatter.formatter.Format(entry)
	return msg, err
}
