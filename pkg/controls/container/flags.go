package container

import (
	"fmt"
	"strings"

	"os"

	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/ninedraft/boxofstuff/str"
)

type Flags struct {
	Image     string
	Memory    uint
	CPU       uint
	Env       []string
	Volume    []string
	Configmap []string
}

func (flags Flags) Container() (container.Container, error) {
	var errs []string
	var cont = container.Container{}

	if volumes, err := flags.volumes(); err == nil {
		cont.VolumeMounts = volumes
	} else {
		errs = append(errs, err.Error())
	}

	if configs, err := flags.configs(); err == nil {
		cont.ConfigMaps = configs
	} else {
		errs = append(errs, err.Error())
	}

	if envs, err := flags.envs(); err == nil {
		cont.Env = envs
	} else {
		errs = append(errs, err.Error())
	}
	if len(errs) > 0 {
		return container.Container{}, fmt.Errorf("unable to build container:\n%s\n",
			str.Vector(errs).Map(func(str string) string {
				return " + " + str
			}).Join("\n"))
	}
	return cont, nil
}

func (flags Flags) volumes() ([]model.ContainerVolume, error) {
	var volumes = make([]model.ContainerVolume, 0, len(flags.Volume))
	for _, volumeStr := range flags.Volume {
		var vol, err = parseContainerVolume(volumeStr, "/mnt")
		if err != nil {
			return nil, fmt.Errorf("invalid volume flag: %v", err)
		}
		volumes = append(volumes, vol)
	}
	return volumes, nil
}

func (flags Flags) configs() ([]model.ContainerVolume, error) {
	var configs = make([]model.ContainerVolume, 0, len(flags.Volume))
	for _, volumeStr := range flags.Configmap {
		var vol, err = parseContainerVolume(volumeStr, "/etc")
		if err != nil {
			return nil, fmt.Errorf("invalid configmap flag: %v", err)
		}
		configs = append(configs, vol)
	}
	return configs, nil
}

func (flags Flags) envs() ([]model.Env, error) {
	var envs = make([]model.Env, 0, len(flags.Env))
	for _, envString := range flags.Env {
		var tokens = str.SplitS(envString, ":", 2).
			Map(strings.TrimSpace).
			Map(os.ExpandEnv)
		if tokens.Len() != 2 {
			return nil, fmt.Errorf("invalid env flag %q: expected NAME:VALUE notation", envString)
		}
		envs = append(envs, model.Env{
			Name:  tokens[0],
			Value: tokens[1],
		})
	}
	return envs, nil
}

// NAME:PATH
// NAME
func parseContainerVolume(volumeStr, defaultMountPath string) (model.ContainerVolume, error) {
	var tokens = str.SplitS(volumeStr, ":", 2).Map(strings.TrimSpace)
	switch tokens.Len() {
	case 1: // NAME
		return model.ContainerVolume{
			Name:      tokens[0],
			MountPath: defaultMountPath + "/" + tokens[0],
		}, nil
	case 2: // NAME:PATH
		return model.ContainerVolume{
			Name:      tokens[0],
			MountPath: tokens[1],
		}, nil
	default:
		return model.ContainerVolume{}, fmt.Errorf("unable to parse %q", volumeStr)
	}
}
