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
	var repeat bool
	for i := uint(0); i == 0 || i < maxTimes; i++ {
		repeat, err = fn()
		if repeat {
			waitNextAttempt(i)
			continue
		}
		return err
	}
	return err
}
