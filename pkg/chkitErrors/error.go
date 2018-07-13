package chkitErrors

import (
	"bytes"
	"fmt"

	"strings"
)

type Err string

var (
	_ error      = Err("")
	_ ErrMatcher = Err("")

	_ error      = &Wrapper{}
	_ ErrMatcher = &Wrapper{}
)

type ErrMatcher interface {
	error
	Match(...error) bool
}

func (err Err) Error() string {
	return string(err)
}
func (err Err) Errors() []error {
	return []error{err}
}

func (err Err) ExitCode() int {
	return 1
}
func (err Err) Wrap(errs ...error) *Wrapper {
	return Wrap(err, errs...)
}

func (err Err) Comment(comments ...string) *Wrapper {
	return Wrap(err).Comment(comments...)
}

func (err Err) CommentF(f string, args ...interface{}) *Wrapper {
	return Wrap(err).Comment(fmt.Sprintf(f, args...))
}

func (err Err) Match(errs ...error) bool {
	for _, er := range errs {
		switch er := er.(type) {
		case *Wrapper:
			if er.main == err {
				return true
			}
		default:
			if err == er {
				return true
			}
		}
	}
	return false
}

type Wrapper struct {
	main          error
	reasons       []error
	comments      []string
	cachedMessage string
}

func (wrapper *Wrapper) Comment(comments ...string) *Wrapper {
	wrapper.comments = append(wrapper.comments, comments...)
	return wrapper
}

func (wrapper *Wrapper) CommentF(f string, args ...interface{}) *Wrapper {
	return wrapper.Comment(fmt.Sprintf(f, args...))
}

func Wrap(err error, reasons ...error) *Wrapper {
	return &Wrapper{
		main:    err,
		reasons: reasons,
	}
}

func (wrapper *Wrapper) AddReasons(reasons ...error) *Wrapper {
	wrapper.reasons = append(wrapper.reasons, reasons...)
	return wrapper
}

func (wrapper *Wrapper) AddReasonF(f string, vals ...interface{}) *Wrapper {
	return wrapper.AddReasons(fmt.Errorf(f, vals...))
}
func (wrapper *Wrapper) Error() string {
	if wrapper.cachedMessage != "" {
		return wrapper.cachedMessage
	}
	buf := bytes.NewBufferString(wrapper.main.Error())

	if len(wrapper.comments) > 0 {
		buf.WriteString(", ")
	}
	buf.WriteString(strings.Join(wrapper.comments, ", "))

	if len(wrapper.reasons) > 0 {
		buf.WriteString(": ")
	}
	for i, reason := range wrapper.reasons {
		if i != 0 {
			buf.WriteString(", " + reason.Error())
		} else {
			buf.WriteString(reason.Error())
		}
	}
	wrapper.cachedMessage = buf.String()
	return wrapper.cachedMessage
}

func (wrapper *Wrapper) Match(errs ...error) bool {
	for _, err := range errs {
		if wrapper.main == err {
			return true
		}
	}
	return false
}

func (wrapper *Wrapper) Errors() []error {
	errs := make([]error, len(wrapper.reasons))
	copy(errs, wrapper.reasons)
	return errs
}

func (wrapper *Wrapper) ExitCode() int {
	return 1
}
