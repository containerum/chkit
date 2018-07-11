package straw

import (
	"fmt"

	"errors"
)

func Catch(fun func()) (err error) {
	defer func() {
		switch er := recover().(type) {
		case nil:
			return
		case error:
			err = er
			return
		case string:
			err = errors.New(er)
		default:
			err = fmt.Errorf("%v", er)
		}
	}()
	fun()
	return nil
}

func CatchErr(fun func() error) (err error) {
	defer func() {
		switch er := recover().(type) {
		case nil:
			return
		case error:
			err = er
			return
		case string:
			err = errors.New(er)
		default:
			err = fmt.Errorf("%v", er)
		}
	}()
	return fun()
}
