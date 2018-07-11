package chkitErrors

import (
	"fmt"
)

type Fatality interface {
	error
	Unwrap() error
	Fatal() string
	Map(func(err error) error) Fatality
}

type fatality struct {
	error
}

func (fatal fatality) Unwrap() error {
	return fatal.error
}

func (fatal fatality) Map(op func(err error) error) Fatality {
	if fatal.error != nil {
		return fatality{op(fatal)}
	}
	return fatality{}
}

func (fatal fatality) Fatal() string {
	return fatal.Error()
}

func Fatal(err error) Fatality {
	return fatality{err}
}

func FatalString(err string, args ...interface{}) Fatality {
	return fatality{fmt.Errorf(err, args...)}
}
