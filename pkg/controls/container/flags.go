package container

import (
	"fmt"
	"os"
	"strings"

	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/ninedraft/boxofstuff/str"
)

type Flags struct {
	Image     string   `desc:"container image"`
	Memory    uint     `desc:"container memory limit, Mb"`
	CPU       uint     `desc:"container CPU limit, mCPU"`
	Env       []string `desc:"container environment variables, NAME:VALUE, 'NAME:$HOST_ENV' or '$HOST_ENV' (to user host env).\nWARNING: single quotes are required to prevent env from interpolation"`
	Commands  []string `desc:"container commands,\nCONTAINER_NAME@VALUE in case of multiple containers or VALUE in case of one container"`
	Volume    []string `desc:"container volume mounts, VOLUME:MOUNT_PATH or VOLUME (then MOUNT_PATH is /mnt/VOLUME)"`
	Configmap []string `desc:"container configmap mount, CONFIG:MOUNT_PATH or CONFIG (then MOUNTPATH is /etc/CONFIG)"`
}

func (flags Flags) Container() (container.Container, error) {
	var errs []error
	var cont = container.Container{}

	cont.Limits = flags.limits()
	cont.Image = flags.Image

	if volumes, err := flags.volumes(); err == nil {
		cont.VolumeMounts = volumes
	} else {
		errs = append(errs, err)
	}

	if configs, err := flags.configs(); err == nil {
		cont.ConfigMaps = configs
	} else {
		errs = append(errs, err)
	}

	if envs, err := flags.envs(); err == nil {
		cont.Env = envs
	} else {
		errs = append(errs, err)
	}

	if cmds, err := flags.cmds(); err == nil {
		cont.Commands = cmds
	} else {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return container.Container{}, fmt.Errorf("unable to build container:\n%s\n",
			str.FromErrs(errs...).Map(str.Prefix(" + ")).Join("\n"))
	}
	return cont, nil
}

func (flags Flags) limits() model.Resource {
	return model.Resource{
		CPU:    flags.CPU,
		Memory: flags.Memory,
	}
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
		var env, err = parseContainerEnv(envString)
		if err != nil {
			return nil, err
		}
		envs = append(envs, env)
	}
	return envs, nil
}

func (flags Flags) cmds() ([]string, error) {
	var cmds = make([]string, 0, len(flags.Commands))
	for _, cmdString := range flags.Commands {
		cmds = append(cmds, cmdString)
	}
	return cmds, nil
}

func parseContainerEnv(envStr string) (model.Env, error) {
	var tokens = str.SplitS(envStr, ":", 2).
		Map(strings.TrimSpace)

	switch tokens.Len() {
	case 1:
		var name = strings.TrimPrefix(tokens[0], "$")
		return model.Env{
			Name:  name,
			Value: os.Getenv(name),
		}, nil
	case 2:
		return model.Env{
			Name:  tokens[0],
			Value: os.ExpandEnv(tokens[1]),
		}, nil
	default:
		return model.Env{}, fmt.Errorf("invalid env %q: expected NAME:VALUE, NAME:$HOST_ENV or $HOST_ENV", envStr)
	}
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
		return model.ContainerVolume{}, fmt.Errorf("invalid volume %q: expected NAME:PATH or NAME", volumeStr)
	}
}
