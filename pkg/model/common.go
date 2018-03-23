package model

import (
	"time"
)

const (
	CreationTimeFormat = time.RFC822
)

func TimestampFormat(timestamp time.Time) string {
	return timestamp.UTC().Format(CreationTimeFormat)
}
