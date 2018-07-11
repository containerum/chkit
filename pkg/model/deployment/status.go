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

func (status *Status) ToKube() *kubeModel.DeploymentStatus {
	if status == nil {
		return nil
	}
	return &kubeModel.DeploymentStatus{
		Replicas:            int(status.NonTerminated),
		ReadyReplicas:       int(status.ReadyReplicas),
		AvailableReplicas:   int(status.AvailableReplicas),
		UnavailableReplicas: int(status.UnavailableReplicas),
		UpdatedReplicas:     int(status.UpdatedReplicas),
	}
}
