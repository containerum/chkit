package err

import (
	"bytes"
	"fmt"
)

type Err string

func (err Err) Error() string {
	return string(err)
}

// Wrap method adds context errors and return Wrapper
func (err Err) Wrap(errs ...error) *Wrapper {
	return &Wrapper{
		origin: err,
		chain:  errs,
	}
}

// Wrapf method consumes format string, vals
// and wrap Err with error with provided message
func (err Err) Wrapf(formats string, vals ...interface{}) *Wrapper {
	return err.Wrap(fmt.Errorf(formats, vals...))
}

// Wrapper contain an origin error and context errors
type Wrapper struct {
	origin       error
	cachedErrMsg string
	chain        []error
}

// Error method returns msg, generated from origin error
// and context errors. It caches and reuse the message, so
// it's normal to call Error multiple times
func (wrapper *Wrapper) Error() string {
	const delim = " "
	if wrapper.cachedErrMsg == "" {
		buf := bytes.NewBufferString(wrapper.origin.Error())
		if len(wrapper.chain) > 0 {
			buf.WriteString(":" + delim)
			for _, err := range wrapper.chain {
				buf.WriteString(delim + err.Error())
			}
			wrapper.cachedErrMsg = buf.String()
		}
	}
	return wrapper.cachedErrMsg
}

func (wrapper Wrapper) Origin() error {
	return wrapper.origin
}

func (wrapper Wrapper) ErrChain() []error {
	return wrapper.chain
}
