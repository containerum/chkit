package compose

import "time"

type Deploy struct {
	Replicas uint
}

type UpdateDeploy struct {
	Deploy
	Delay time.Duration
}
