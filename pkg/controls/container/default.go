package container

import (
	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/ninedraft/boxofstuff/str"
)

func Default(cont container.Container) container.Container {
	cont = cont.Copy()
	if cont.Name == "" {
		cont.Name = str.Vector{namegen.Color(), cont.ImageName(), namegen.Aster()}.
			Filter(str.Longer(0))[:2].
			Join("-")
	}
	if cont.Limits.Memory == 0 {
		cont.Limits.Memory = 256 // Mb
	}
	if cont.Limits.CPU == 0 {
		cont.Limits.CPU = 200 // mCPU
	}
	for i, vol := range cont.VolumeMounts {
		if vol.MountPath == "" && vol.Name != "" {
			vol.MountPath = "/mnt/" + vol.Name
		}
		cont.VolumeMounts[i] = vol
	}
	for i, config := range cont.ConfigMaps {
		if config.MountPath == "" && config.Name != "" {
			config.MountPath = "/etc/" + config.Name
		}
		cont.ConfigMaps[i] = config
	}
	for i, port := range cont.Ports {
		if port.Protocol == "" {
			port.Protocol = model.TCP
		}
		if port.Port <= 0 {
			port.Port = 80
		}
		if port.Name == "" {
			port.Name = namegen.Aster() + "-" + namegen.Physicist()
		}
		cont.Ports[i] = port
	}
	return cont
}
