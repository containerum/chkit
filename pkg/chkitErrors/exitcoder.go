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

func NewExitCoder(err error) cli.ExitCoder {
	switch err := err.(type) {
	case cli.ExitCoder:
		return err
	default:
		return &ExitCoder{
			Code: 1,
			Err:  err,
		}
	}
}
func (coder *ExitCoder) Error() string {
	return coder.Err.Error()
}

func (coder *ExitCoder) ExitCode() int {
	return coder.Code
}
