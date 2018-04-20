package deployment

import (
	"time"

	"github.com/sirupsen/logrus"

	kubeModel "git.containerum.net/ch/kube-client/pkg/model"
)

type Status struct {
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Replicas            uint
	ReadyReplicas       uint
	AvailableReplicas   uint
	UnavailableReplicas uint
	UpdatedReplicas     uint
}

func StatusFromKubeStatus(kubeStatus kubeModel.DeploymentStatus) Status {
	createdAt, err := time.Parse(time.RFC3339, kubeStatus.CreatedAt)
	if err != nil {
		logrus.WithError(err).Debugf("invalid created_at timestamp")
		createdAt = time.Unix(0, 0)
	}
	updatedAt, err := time.Parse(time.RFC3339, kubeStatus.UpdatedAt)
	if err != nil {
		logrus.WithError(err).Debugf("invalid updated_at timestamp")
		updatedAt = time.Unix(0, 0)
	}
	return Status{
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
		Replicas:            uint(kubeStatus.Replicas),
		AvailableReplicas:   uint(kubeStatus.AvailableReplicas),
		UnavailableReplicas: uint(kubeStatus.UnavailableReplicas),
		UpdatedReplicas:     uint(kubeStatus.UpdatedReplicas),
	}
}
