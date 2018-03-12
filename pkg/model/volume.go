package model

import (
	"time"
)

type Volume struct {
	Label     string
	CreatedAt time.Time
	Access    string
	Replicas  uint
	Storage   uint
}
