package model

import (
	"fmt"
	"time"
)

const (
	TimestampFormat = time.RFC3339
	Indent          = "  "
)

func Age(timestamp time.Time) string {
	if timestamp.Equal(time.Unix(0, 0)) {
		return "unknown"
	}
	age := time.Now().Sub(timestamp)
	var ageString string
	const year = 365 * 24
	switch {
	case age.Hours() > year:
		years := uint64(age.Hours()) / year
		ageString = fmt.Sprintf("%dy", years)
	case age.Hours() > 24 && age.Hours() < year:
		days := uint64(age.Hours()) / 24
		ageString = fmt.Sprintf("%dd", days)
	case age.Hours() <= 24 && age.Hours() > 1:
		hours := uint64(age.Hours())
		ageString = fmt.Sprintf("%dh", hours)
	case age.Hours() < 1 && age.Minutes() > 1:
		minutes := uint64(age.Minutes())
		ageString = fmt.Sprintf("%dm", minutes)
	default:
		seconds := uint64(age.Seconds())
		ageString = fmt.Sprintf("%ds", seconds)
	}
	return ageString
}
