package deployment

import (
	kubeModel "github.com/containerum/kube-client/pkg/model"
)

type Status struct {
	NonTerminated       uint
	ReadyReplicas       uint
	AvailableReplicas   uint
	UnavailableReplicas uint
	UpdatedReplicas     uint
}

func StatusFromKubeStatus(kubeStatus kubeModel.DeploymentStatus) Status {
	return Status{
		NonTerminated:       uint(kubeStatus.Replicas),
		AvailableReplicas:   uint(kubeStatus.AvailableReplicas),
		UnavailableReplicas: uint(kubeStatus.UnavailableReplicas),
		UpdatedReplicas:     uint(kubeStatus.UpdatedReplicas),
	}
}
