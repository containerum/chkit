package deployment

import (
	"fmt"
	"strings"
	"time"

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
	return Status{
		CreatedAt:           time.Unix(kubeStatus.CreatedAt, 0),
		UpdatedAt:           time.Unix(kubeStatus.UpdatedAt, 0),
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
		"Ready: " + fmt.Sprintf("%d/%d", status.ReadyReplicas, status.Replicas),
		"Updated: " + fmt.Sprintf("%d/%d", status.UpdatedReplicas, status.Replicas),
	}, "\n")
}

func (status *Status) ColumnWhen() string {
	if status == nil {
		return "unknown"
	}
	return strings.Join([]string{
		"Created: " + status.CreatedAt.Format(model.CreationTimeFormat),
		"Updated: " + status.UpdatedAt.Format(model.CreationTimeFormat),
	}, "\n")
}
