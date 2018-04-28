package pod

import (
	"fmt"
	"time"

	kubeModel "github.com/containerum/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/model"
)

type Pod struct {
	Name       string
	Hostname   string
	Containers []string
	Status     Status
	CreatedAt  time.Time
	origin     kubeModel.Pod
}

func PodFromKube(pod kubeModel.Pod) Pod {
	var containers []string
	for _, container := range pod.Containers {
		containers = append(containers,
			fmt.Sprintf("%s [%s]",
				container.Name,
				container.Image))
	}
	hostname := ""
	if pod.Hostname != nil {
		hostname = *pod.Hostname
	}
	var status Status
	if pod.Status != nil {
		status = StatusFromKube(*pod.Status)
	}
	var createdAt time.Time
	if pod.CreatedAt != nil {
		createdAt, _ = time.Parse(model.TimestampFormat, *pod.CreatedAt)
	}
	return Pod{
		Name:       pod.Name,
		Hostname:   hostname,
		Containers: containers,
		Status:     status,
		CreatedAt:  createdAt,
		origin:     pod,
	}
}
