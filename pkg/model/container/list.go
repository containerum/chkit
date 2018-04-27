package container

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
