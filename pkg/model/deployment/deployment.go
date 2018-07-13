package deployment

import (
	"fmt"
	"time"

	"github.com/blang/semver"
	model2 "github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/ninedraft/boxofstuff/str"
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
	SolutionID  string
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
		SolutionID:  kubeDeployment.SolutionID,
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
		SolutionID:  depl.SolutionID,
	}
	return kubeDepl
}

func (depl *Deployment) StatusString() string {
	var status string
	if depl.Active {
		if depl.Status != nil {
			status += fmt.Sprintf("Running pods: %2d/%d\n"+
				"CPU usage:   %4d mCPU\n"+
				"RAM usage:   %4d Mb",
				depl.Status.AvailableReplicas, depl.Replicas,
				depl.TotalCPU,
				depl.TotalMemory)
		} else {
			status += fmt.Sprintf("local\nreplicas %d", depl.Replicas)
		}
	}
	if len(depl.Containers) == 0 {
		status += "\n!MISSING CONTAINERS!"
	}
	return str.Vector{status, "inactive"}.FirstNonEmpty()
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
