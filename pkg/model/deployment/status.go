package deployment

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	kubeModel "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/model"
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

func (status *Status) ColumnReplicas() string {
	if status == nil {
		return "unknown"
	}
	return strings.Join([]string{
		"Available: " + fmt.Sprintf("%d/%d", status.AvailableReplicas, status.Replicas),
	}, "\n")
}

func (status *Status) ColumnWhen() string {
	if status == nil {
		return "unknown"
	}
	return strings.Join([]string{
		"Created: " + model.Age(status.CreatedAt),
		"Updated: " + model.Age(status.UpdatedAt),
	}, "\n")
}
