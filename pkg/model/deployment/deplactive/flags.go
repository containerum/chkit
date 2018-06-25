package deplactive

import (
	"fmt"
	"strconv"
	"strings"

	chkitContainer "github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/ninedraft/boxofstuff/str"
)

func missingContainerNameErr(flag, value string) error {
	return fmt.Errorf("invalid %s flag value %q: container name can be omitted only if only one container defined", flag, value)
}

func invalidFlagValueErr(flag, value, reason string) error {
	return fmt.Errorf("invalid %s flag value %q: %s", flag, value, reason)
}

type Flags struct {
	Force bool   `flag:"force f" desc:"suppress confirmation, optional"`
	File  string `desc:"file with configmap data, .json, .yaml, .yml, optional"`
	// Output   string `flag:"output o" desc:"output format, json/yaml"`
	Replicas      uint     `desc:"deployment replicas, optional"` // deployment
	Name          string   `desc:"deployment name, optional"`     // deployment
	ContainerName string   `desc:"container name in case of single container"`
	Image         []string `desc:"container image,\nCONTAINER_NAME@IMAGE in case of multiple containers or IMAGE in case of one container"` // container +

	Env []string `desc:"container environment variable,\nCONTAINER_NAME@KEY:VALUE in case of multiple containers or KEY:VALUE in case of one container"` // container +

	Memory []string `desc:"container memory limit, Mb,\nCONTAINER_NAME@MEMORY in case of multiple containers or MEMORY in case of one container"` // container +

	CPU []string `desc:"container memory limit, mCPU,\nCONTAINER_NAME@CPU in case of multiple containers or CPU in case of one container"` // container +

	Volume []string `desc:"container volume,\nCONTAINER_NAME@VOLUME_NAME@MOUNTPATH in case of multiple containers or\nVOLUME_NAME@MOUNTPATH or VOLUME_NAME in case of one container.\nIf MOUNTPATH is omitted, then use /mnt/VOLUME_NAME as mountpath"` // container +

	Configmap []string `desc:"container configmap, CONTAINER_NAME@CONFIGMAP_NAME@MOUNTPATH in case of multiple containers or\nCONFIGMAP_NAME@MOUNTPATH or CONFIGMAP_NAME in case of one container.\nIf MOUNTPATH is omitted, then use /etc/CONFIGMAP_NAME as mountpath"` // container +

	containers map[string]chkitContainer.Container
}

func FlagsFromDeployment(oldflags Flags, depl deployment.Deployment) Flags {
	var containers = make(map[string]chkitContainer.Container, len(depl.Containers))
	for _, container := range depl.Containers {
		containers[container.Name] = container
	}
	ret := oldflags
	ret.Name = depl.Name
	ret.Replicas = uint(depl.Replicas)
	ret.containers = containers
	return ret
}

func (flags Flags) Deployment() (deployment.Deployment, error) {
	var containers, err = flags.BuildContainers()
	if err != nil {
		return deployment.Deployment{}, err
	}
	return deployment.Deployment{
		Name:       flags.Name,
		Replicas:   int(flags.Replicas),
		Containers: containers,
	}, nil
}

func (flags Flags) BuildContainers() (chkitContainer.ContainerList, error) {
	if flags.containers == nil {
		flags.containers = map[string]chkitContainer.Container{}
	}
	if err := flags.extractImages(); err != nil {
		return nil, err
	}
	if err := flags.extractMemory(); err != nil {
		return nil, err
	}
	if err := flags.extractCPU(); err != nil {
		return nil, err
	}
	if err := flags.extractEnvs(); err != nil {
		return nil, err
	}
	if err := flags.extractVolumes(); err != nil {
		return nil, err
	}
	if err := flags.extractConfigmaps(); err != nil {
		return nil, err
	}

	var list = make(chkitContainer.ContainerList, 0, len(flags.containers))
	for containerName, container := range flags.containers {
		containerName = str.Vector{
			containerName,
			flags.ContainerName,
			namegen.Color() + "-" + container.Image}.FirstNonEmpty()
		container.Name = containerName
		list = append(list, container)
	}
	return list, nil
}

func (flags Flags) extractCPU() error {
	for _, cpuValue := range flags.CPU {
		var container, cpuStr = flags.extractContainerAndValue(cpuValue)
		//	if container == "" {
		//		return missingContainerNameErr("--cpu", cpuStr)
		//	}
		var cont = flags.containers[container]
		var cpu, err = strconv.ParseUint(cpuStr, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid --cpu flag value: %v", err)
		}
		cont.Limits.CPU = uint(cpu)
		flags.containers[container] = cont
	}
	return nil
}

func (flags Flags) extractImages() error {
	for _, image := range flags.Image {
		var container, imageName = flags.extractContainerAndValue(image)
		//	if container == "" {
		//		return missingContainerNameErr("--image", imageName)
		//	}
		var cont = flags.containers[container]
		cont.Image = imageName
		flags.containers[container] = cont
	}
	return nil
}

func (flags Flags) extractMemory() error {
	for _, memValue := range flags.Memory {
		var container, memStr = flags.extractContainerAndValue(memValue)
		// if container == "" {
		//	return missingContainerNameErr("--memory", memValue)
		// }
		var cont = flags.containers[container]
		var mem, err = strconv.ParseUint(memStr, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid --memory flag value: %v", err)
		}
		cont.Limits.Memory = uint(mem)
		flags.containers[container] = cont
	}
	return nil
}

func (flags Flags) extractEnvs() error {
	for _, envValue := range flags.Env {
		var container, envPair = flags.extractContainerAndValue(envValue)
		// if container == "" {
		//	return missingContainerNameErr("--env", envValue)
		// }
		var env, err = parseEnv(envPair)
		if err != nil {
			return err
		}
		var cont = flags.containers[container]
		cont.AddEnv(env)
		flags.containers[container] = cont
	}
	return nil
}

func (flags Flags) extractVolumes() error {
	for _, volumeValue := range flags.Volume {
		var container, volumeName, mountPath = flags.extractContainerAndVolumeNameAndMountPath(volumeValue)
		//	if container == "" {
		//		return missingContainerNameErr("--volume", volumeValue)
		//	}
		if volumeName == "" {
			return invalidFlagValueErr(
				"--volume",
				volumeValue,
				"expected [CONTAINER]@[VOLUME_NAME]@[VOLUME_MOUNT] or [CONTAINER]@[VOLUME_NAME]")
		}
		if mountPath == "" {
			mountPath = "/mnt/" + volumeName
		}
		var cont = flags.containers[container]
		cont.VolumeMounts = append(cont.VolumeMounts, model.ContainerVolume{
			Name:      volumeName,
			MountPath: mountPath,
		})
		flags.containers[container] = cont
	}
	return nil
}

func (flags Flags) extractConfigmaps() error {
	for _, configmapValue := range flags.Configmap {
		var container, configName, mountPath = flags.extractContainerAndVolumeNameAndMountPath(configmapValue)
		//	if container == "" {
		//		return missingContainerNameErr("--configmap", configmapValue)
		//	}
		if configName == "" {
			return invalidFlagValueErr(
				"--volume",
				configmapValue,
				"expected [CONTAINER]@[CONFIGMAP_NAME]@[CONFIGMAP_MOUNT] or [CONTAINER]@[CONFIGMAP_NAME]")
		}
		if mountPath == "" {
			mountPath = "/etc/" + configName
		}
		var cont = flags.containers[container]
		cont.ConfigMaps = append(cont.ConfigMaps, model.ContainerVolume{
			Name:      configName,
			MountPath: mountPath,
		})
		flags.containers[container] = cont
	}
	return nil
}

func (flags Flags) extractContainerAndValue(str string) (container, value string) {
	var tokens = strings.SplitN(str, "@", 2)
	if len(tokens) == 2 {
		return tokens[0], tokens[1]
	}
	return flags.ContainerName, str
}

// [CONTAINER]@[VOLUME_NAME]@[VOLUME_MOUNT]
// [CONTAINER]@[VOLUME_NAME]
// [VOLUME_NAME]
func (flags Flags) extractContainerAndVolumeNameAndMountPath(str string) (container, volume, mount string) {
	var tokens = strings.SplitN(str, "@", 3)
	switch len(tokens) {
	case 3: //  [CONTAINER]@[VOLUME_NAME]@[VOLUME_MOUNT]
		return tokens[0], tokens[1], tokens[2]
	case 2: // [CONTAINER]@[VOLUME_NAME]
		return tokens[0], tokens[1], ""
	case 1: // [VOLUME_NAME]
		return flags.ContainerName, str, ""
	default:
		return "", "", ""
	}
}

func parseEnv(envString string) (model.Env, error) {
	var tokens = strings.SplitN(envString, ":", 2)
	if len(tokens) != 2 {
		return model.Env{}, fmt.Errorf("invalid env pair: expect $KEY:$VALUE, got %q", envString)
	}
	var key, value = tokens[0], tokens[1]
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	return model.Env{
		Name:  key,
		Value: value,
	}, nil
}
