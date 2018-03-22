package deployment

import (
	"fmt"
	"strings"
	"time"

	"github.com/containerum/chkit/pkg/model"
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

func (status *Status) String() string {
	if status == nil {
		return "unknown"
	}
	return strings.Join([]string{
		"Created: " + status.CreatedAt.Format(model.CreationTimeFormat),
		"Updated: " + status.UpdatedAt.Format(model.CreationTimeFormat),
		"Available replicas: " + fmt.Sprintf("%d/%d", status.AvailableReplicas, status.Replicas),
		"Ready replicas: " + fmt.Sprintf("%d/%d", status.ReadyReplicas, status.Replicas),
	}, "\n")
}
