package deployment

import (
	"fmt"
	"time"

	"github.com/blang/semver"
	model2 "github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/kube-client/pkg/model"
)

type Deployment struct {
	Name        string
	Replicas    int
	Status      *Status
	Active      bool
	Version     semver.Version
	CreatedAt   time.Time
	TotalCPU    uint
	TotalMemory uint
	Containers  container.ContainerList
}

func DeploymentFromKube(kubeDeployment model.Deployment) Deployment {
	var status *Status
	if kubeDeployment.Status != nil {
		st := StatusFromKubeStatus(*kubeDeployment.Status)
		status = &st
	}
	var timestamp time.Time
	if t, err := time.Parse(time.RFC3339, kubeDeployment.CreatedAt); err == nil {
		timestamp = t
	}
	containers := make([]container.Container, 0, len(kubeDeployment.Containers))
	for _, kubeContainer := range kubeDeployment.Containers {
		containers = append(containers, container.Container{Container: kubeContainer}.Copy())
	}
	return Deployment{
		Name:        kubeDeployment.Name,
		Replicas:    kubeDeployment.Replicas,
		Status:      status,
		Containers:  containers,
		Version:     kubeDeployment.Version,
		Active:      kubeDeployment.Active,
		CreatedAt:   timestamp,
		TotalCPU:    kubeDeployment.TotalCPU,
		TotalMemory: kubeDeployment.TotalMemory,
	}
}

func (depl *Deployment) ToKube() model.Deployment {
	containers := make([]model.Container, 0, len(depl.Containers))
	for _, cont := range depl.Containers {
		containers = append(containers, cont.Container)
	}
	kubeDepl := model.Deployment{
		Name:        depl.Name,
		Replicas:    int(depl.Replicas),
		Containers:  containers,
		Version:     depl.Version,
		Active:      depl.Active,
		TotalMemory: depl.TotalMemory,
		TotalCPU:    depl.TotalCPU,
		CreatedAt:   depl.CreatedAt.Format(model2.TimestampFormat),
		Status:      depl.Status.ToKube(),
	}
	return kubeDepl
}

func (depl *Deployment) StatusString() string {
	if depl.Active {
		if depl.Status != nil {
			return fmt.Sprintf("running: %2d/%d\n"+
				"CPU    : %4d mCPU\n"+
				"MEMORY : %4d Mb",
				depl.Status.NonTerminated, depl.Replicas,
				depl.TotalCPU,
				depl.TotalMemory)
		} else {
			return fmt.Sprintf("local\nreplicas %d", depl.Replicas)
		}
	}
	return "inactive"
}

func (depl Deployment) Copy() Deployment {
	var status *Status
	if depl.Status != nil {
		var s = *depl.Status
		status = &s
	}
	var version = depl.Version
	version.Build = append([]string{}, version.Build...)
	version.Pre = append([]semver.PRVersion{}, version.Pre...)
	return Deployment{
		Name:        depl.Name,
		TotalCPU:    depl.TotalCPU,
		TotalMemory: depl.TotalMemory,
		CreatedAt:   depl.CreatedAt,
		Replicas:    depl.Replicas,
		Active:      depl.Active,
		Status:      status,
		Version:     version,
		Containers:  depl.Containers.Copy(),
	}
}
