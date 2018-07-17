package container

import (
	"github.com/containerum/chkit/pkg/model/container"
)

type ReplaceFlags struct {
	Flags
	DeleteConfigmap []string `desc:"configmap to delete"`
	DeleteVolume    []string `desc:"volume to delete"`
	DeleteEnv       []string `desc:"environment to delete"`
}

func (flags ReplaceFlags) Patch(cont container.Container) (container.Container, error) {
	var left, err = flags.Container()
	if err != nil {
		return container.Container{}, err
	}
	var patched = left.Patch(cont)
	{
		var configs = patched.ConfigMountsMap()
		for _, config := range flags.DeleteConfigmap {
			delete(configs, config)
		}
		patched.ConfigMaps = patched.ConfigMaps[:0]
		for _, config := range configs {
			patched.ConfigMaps = append(patched.ConfigMaps, config)
		}
	}
	{
		var volumes = patched.VolumeMountsMap()
		for _, volume := range flags.DeleteVolume {
			delete(volumes, volume)
		}
		patched.VolumeMounts = patched.VolumeMounts[:0]
		for _, volume := range volumes {
			patched.VolumeMounts = append(patched.VolumeMounts, volume)
		}
	}
	{
		var envs = patched.GetEnvMap()
		for _, env := range flags.DeleteEnv {
			delete(envs, env)
		}
		patched.PutEnvMap(envs)
	}
	return patched, nil
}
