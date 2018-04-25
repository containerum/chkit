package deployment

import (
	"git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/model/container"
)

type Deployment struct {
	Name       string
	Replicas   int
	Status     *Status
	Containers []container.Container
	origin     *model.Deployment
}

func DeploymentFromKube(kubeDeployment model.Deployment) Deployment {
	var status *Status
	if kubeDeployment.Status != nil {
		st := StatusFromKubeStatus(*kubeDeployment.Status)
		status = &st
	}
	containers := make([]container.Container, 0, len(kubeDeployment.Containers))
	for _, kubeContainer := range kubeDeployment.Containers {
		containers = append(containers, container.Container{Container: kubeContainer})
	}
	return Deployment{
		Name:       kubeDeployment.Name,
		Replicas:   kubeDeployment.Replicas,
		Status:     status,
		Containers: containers,
		origin:     &kubeDeployment,
	}
}

func (depl *Deployment) ToKube() model.Deployment {
	containers := make([]model.Container, 0, len(depl.Containers))
	for _, cont := range depl.Containers {
		containers = append(containers, cont.Container)
	}
	kubeDepl := model.Deployment{
		Name:       depl.Name,
		Replicas:   int(depl.Replicas),
		Containers: containers,
	}
	depl.origin = &kubeDepl
	return kubeDepl
}
