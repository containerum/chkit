package chkitErrors

import (
	"gopkg.in/urfave/cli.v2"
)

var (
	_ cli.ExitCoder = &ExitCoder{}
)

type ExitCoder struct {
	Err  error
	Code int
}

func NewExitCoder(err error) *ExitCoder {
	return &ExitCoder{
		Err:  err,
		Code: 1,
	}
}
func (coder *ExitCoder) Error() string {
	return coder.Err.Error()
}

func (coder *ExitCoder) ExitCode() int {
	return coder.Code
}
