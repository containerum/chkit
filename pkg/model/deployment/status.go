package deployment

import (
	"time"
)

type Status struct {
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Replicas            int
	ReadyReplicas       int
	AvailableReplicas   int
	UnavailableReplicas int
	UpdatedReplicas     int
}
