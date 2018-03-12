package model

import (
	"time"
)

type Namespace struct {
	CreatedAt *time.Time
	Label     string
	Access    string
	Volumes   []Volume
}
