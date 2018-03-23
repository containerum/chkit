package deployment

import (
	"git.containerum.net/ch/kube-client/pkg/model"
)

type Deployment struct {
	Name       string
	Replicas   uint
	Status     *Status
	Containers []Container
	origin     model.Deployment
}

func DeploymentFromKube(kubeDeployment model.Deployment) Deployment {
	var status *Status
	if kubeDeployment.Status != nil {
		st := StatusFromKubeStatus(*kubeDeployment.Status)
		status = &st
	}
	containers := make([]Container, 0, len(kubeDeployment.Containers))
	for _, kubeContainer := range kubeDeployment.Containers {
		containers = append(containers, Container{kubeContainer})
	}
	return Deployment{
		Name:       kubeDeployment.Name,
		Replicas:   uint(kubeDeployment.Replicas),
		Status:     status,
		Containers: containers,
		origin:     kubeDeployment,
	}
}
