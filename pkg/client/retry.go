package chClient

import (
	"time"
)

func waitNextAttempt(attempt uint) {
	duration := 500 * time.Duration(attempt+1) * time.Millisecond
	if duration < time.Minute {
		time.Sleep(duration)
	} else {
		time.Sleep(time.Minute)
	}
}

func retry(maxTimes uint, fn func() (bool, error)) error {
	var err error
	var ok bool
	for i := uint(0); i == 0 || (i < maxTimes && err != nil); i++ {
		ok, err = fn()
		if err == nil || !ok {
			break
		}
		waitNextAttempt(i)
	}
	return err
}
