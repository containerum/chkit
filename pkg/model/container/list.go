package container

import "github.com/containerum/kube-client/pkg/model"

type ContainerList []Container

func (containers ContainerList) Images() []string {
	images := make([]string, 0, len(containers))
	for _, container := range containers {
		images = append(images, container.Image)
	}
	return images
}

func (containers ContainerList) GetByName(name string) (Container, bool) {
	for _, container := range containers {
		if container.Name == name {
			return container, true
		}
	}
	return Container{}, false
}

func (containers ContainerList) Names() []string {
	names := make([]string, 0, len(containers))
	for _, container := range containers {
		names = append(names, container.Name)
	}
	return names
}

func (list ContainerList) ConfigMountsMap() map[string]model.ContainerVolume {
	var mounts = make(map[string]model.ContainerVolume, len(list))
	for _, container := range list {
		for _, config := range container.ConfigMaps {
			mounts[config.Name] = config
		}
	}
	return mounts
}

func (list ContainerList) VolumeMountsMap() map[string]model.ContainerVolume {
	var mounts = make(map[string]model.ContainerVolume, len(list))
	for _, container := range list {
		for _, volume := range container.VolumeMounts {
			mounts[volume.Name] = volume
		}
	}
	return mounts
}

func (list ContainerList) Copy() ContainerList {
	var cp = make(ContainerList, 0, len(list))
	for _, cont := range cp {
		cp = append(cp, cont.Copy())
	}
	return cp
}
