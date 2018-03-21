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

func retry(maxTimes uint, fn func() error) error {
	var err error
	for i := uint(0); i == 0 || (i < maxTimes && err != nil); i++ {
		err = fn()
		if err == nil {
			break
		}
		waitNextAttempt(i)
	}
	return err
}
