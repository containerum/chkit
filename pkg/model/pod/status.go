package pod

import (
	"time"

	"github.com/sirupsen/logrus"

	kubeModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/model"
)

type Status struct {
	Phase        string
	RestartCount int
	StartedAt    time.Time
}

func StatusFromKube(status kubeModels.PodStatus) Status {
	startedAt, err := time.Parse(status.StartAt, model.TimestampFormat)
	if err != nil {
		logrus.WithError(err).Debugf("invalid started_at timestamp")
		startedAt = time.Unix(0, 0)
	}
	return Status{
		Phase:        status.Phase,
		RestartCount: status.RestartCount,
		StartedAt:    startedAt,
	}
}
