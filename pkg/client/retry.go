package chClient

import "time"

func waitNextAttempt(attempt uint) {
	duration := 500 * time.Duration(attempt+1) * time.Millisecond
	if duration < time.Minute {
		time.Sleep(duration)
	} else {
		time.Sleep(time.Minute)
	}
}
