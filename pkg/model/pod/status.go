package pod

import (
	"time"

	"git.containerum.net/ch/kube-client/pkg/model"
)

type Status struct {
	Phase        string
	RestartCount int
	StartedAt    time.Time
}

func StatusFromKube(status model.PodStatus) Status {
	return Status{
		Phase:        status.Phase,
		RestartCount: status.RestartCount,
		StartedAt:    time.Unix(status.StartAt, 0),
	}
}
