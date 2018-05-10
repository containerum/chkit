package pod

import (
	"time"

	"github.com/containerum/chkit/pkg/model"
	kubeModels "github.com/containerum/kube-client/pkg/model"
	"github.com/sirupsen/logrus"
)

type Status struct {
	Phase        string
	RestartCount int
	StartedAt    time.Time
}

func StatusFromKube(status kubeModels.PodStatus) Status {
	startedAt, err := time.Parse(model.TimestampFormat, status.StartAt)
	if err != nil {
		logrus.WithError(err).Debugf("invalid started_at timestamp %q", status.StartAt)
		startedAt = time.Unix(0, 0)
	}
	return Status{
		Phase:        status.Phase,
		RestartCount: status.RestartCount,
		StartedAt:    startedAt,
	}
}
