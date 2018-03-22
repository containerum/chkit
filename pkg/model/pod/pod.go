package pod

import (
	"git.containerum.net/ch/kube-client/pkg/model"
)

type Pod struct {
	Name       string
	Hostname   string
	Containers []string
	Status     Status
}

func PodFromKube(pod model.Pod) Pod {
	containers := []string{}
	for _, container := range pod.Containers {
		containers = append(containers, container.Name)
	}
	hostname := ""
	if pod.Hostname != nil {
		hostname = *pod.Hostname
	}
	var status Status
	if pod.Status != nil {
		status = StatusFromKube(*pod.Status)
	}
	return Pod{
		Name:       pod.Name,
		Hostname:   hostname,
		Containers: containers,
		Status:     status,
	}
}