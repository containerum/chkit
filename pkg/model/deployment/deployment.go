package deployment

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/kube-client/pkg/model"
)

type Deployment struct {
	Name       string
	Replicas   int
	Status     *Status
	Active     bool
	Version    semver.Version
	Containers container.ContainerList
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
		Version:    kubeDeployment.Version,
		Active:     kubeDeployment.Active,
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
		Version:    depl.Version,
		Active:     depl.Active,
	}
	depl.origin = &kubeDepl
	return kubeDepl
}

func (depl *Deployment) StatusString() string {
	if depl.Active {
		if depl.Status != nil {
			return fmt.Sprintf("running %d/%d",
				depl.Status.AvailableReplicas, depl.Replicas)
		} else {
			return fmt.Sprintf("local\nreplicas %d", depl.Replicas)
		}
	}
	return "inactive"
}
