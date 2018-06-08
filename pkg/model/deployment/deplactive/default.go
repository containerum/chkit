package deplactive

import (
	"github.com/blang/semver"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/namegen"
)

func Fill(depl *deployment.Deployment) {
	if depl.Name == "" {
		depl.Name = namegen.Aster() + "-" + namegen.Physicist()
	}
	if depl.Replicas < 1 {
		depl.Replicas = 1
	}
	if depl.Version.Patch == 0 &&
		depl.Version.Minor == 0 &&
		depl.Version.Major == 0 {
		depl.Version = semver.MustParse("1.0.0")
	}
	for i, container := range depl.Containers {
		if container.Name == "" {
			container.Name = namegen.Color() + "-" + namegen.Aster()
		}
		if container.Limits.CPU == 0 {
			container.Limits.CPU = 200
		}
		if container.Limits.Memory == 0 {
			container.Limits.Memory = 256
		}
		for configIndex, config := range container.ConfigMaps {
			if config.MountPath == "" {
				config.MountPath = "/etc/" + config.Name
			}
			container.ConfigMaps[configIndex] = config
		}
		for volumeIndex, volume := range container.VolumeMounts {
			if volume.MountPath == "" {
				volume.MountPath = "/mnt/" + volume.Name
			}
			container.ConfigMaps[volumeIndex] = volume
		}
		depl.Containers[i] = container
	}
}
