package deplactive

import (
	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/model/deployment"
)

func Merge(a, b deployment.Deployment) deployment.Deployment {
	a = a.Copy()
	b = b.Copy()
	if b.Replicas > 0 {
		a.Replicas = b.Replicas
	}
	for _, bContainer := range b.Containers {
		for i, aContainer := range a.Containers {
			if aContainer.Name == bContainer.Name {
				a.Containers[i] = bContainer
				continue
			}
		}
		a.Containers = append(a.Containers, bContainer)
	}
	return a
}

func mergeContainers(a, b container.Container) container.Container {
	a = a.Copy()
	b = b.Copy()
	if b.Image != "" {
		a.Image = b.Image
	}
	if b.Limits.Memory != 0 {
		a.Limits.Memory = b.Limits.Memory
	}
	if b.Limits.CPU != 0 {
		a.Limits.CPU = b.Limits.CPU
	}
	a.Env = b.Env
	a.VolumeMounts = b.VolumeMounts
	a.ConfigMaps = b.ConfigMaps
	return a
}
