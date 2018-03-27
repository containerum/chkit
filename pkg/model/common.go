package model

import (
	"fmt"
	"time"
)

const (
	CreationTimeFormat = time.RFC822
)

func TimestampFormat(timestamp time.Time) string {
	age := time.Now().Sub(timestamp)
	ageString := age.String()
	switch {
	case age.Hours() > 365*24:
		years := uint64(age.Hours()) / (365 * 24)
		days := uint64(age.Hours()) % (365 * 24)
		ageString = fmt.Sprintf("%dy-%03d", years, days)
	case age.Hours() > 24:
		days := uint64(age.Hours()) / 24
		hours := uint64(age.Hours()) % 24
		ageString = fmt.Sprintf("%dd-%02dh", days, hours)
	case age.Hours() <= 24 && age.Hours() > 1:
		hours := uint64(age.Hours())
		minutes := uint64(age.Minutes()) % 60
		ageString = fmt.Sprintf("%dh-%02dm", hours, minutes)
	case age.Hours() < 1 && age.Minutes() > 1:
		minutes := uint64(age.Minutes())
		seconds := uint64(age.Seconds()) % 60
		ageString = fmt.Sprintf("%dm-%02ds", minutes, seconds)
	default:
		seconds := uint64(age.Seconds())
		ageString = fmt.Sprintf("%ds", seconds)
	}
	return ageString
}
