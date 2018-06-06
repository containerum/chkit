package container

import (
	"fmt"
	"strings"

	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type Container struct {
	kubeModels.Container
}

func (container Container) String() string {
	return strings.Join([]string{
		container.Name,
		fmt.Sprintf("CPU %d", container.Limits.CPU),
		fmt.Sprintf("MEM %d", container.Limits.Memory),
	}, " ") + "; "
}

func (container Container) ConfigmapNames() []string {
	var names = make([]string, 0, len(container.ConfigMaps))
	for _, cm := range container.ConfigMaps {
		names = append(names, cm.Name)
	}
	return names
}

func (container Container) ConfigMountsMap() map[string]kubeModels.ContainerVolume {
	var mounts = make(map[string]kubeModels.ContainerVolume, len(container.ConfigMaps))
	for _, config := range container.ConfigMaps {
		mounts[config.MountPath] = config
	}
	return mounts
}

func (container Container) VolumeMountsMap() map[string]kubeModels.ContainerVolume {
	var mounts = make(map[string]kubeModels.ContainerVolume, len(container.VolumeMounts))
	for _, volume := range container.VolumeMounts {
		mounts[volume.MountPath] = volume
	}
	return mounts
}
