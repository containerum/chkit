package container

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/blang/semver"
	"github.com/containerum/chkit/pkg/model"
	kubeModels "github.com/containerum/kube-client/pkg/model"
	"github.com/ninedraft/boxofstuff/str"
)

var (
	_ model.Renderer = Container{}
)

type Container struct {
	kubeModels.Container
}

func (container Container) ToKube() kubeModels.Container {
	return kubeModels.Container(container.Copy().Container)
}

func ImageName(image string) string {
	var img, err = kubeModels.ImageFromString(image)
	if err != nil {
		return image
	}
	var alphaNumerical = regexp.MustCompile("^[a-z0-9]+$")
	return str.SplitS(img.Name, "/", 3).Tail(1).Filter(func(str string) bool {
		return alphaNumerical.MatchString(str)
	}).FirstNonEmpty()
}

func (container Container) ImageName() string {
	return ImageName(container.Image)
}

func (container Container) SemanticVersion() (version semver.Version, ok bool) {
	var vStr = container.Version()
	if version, err := semver.ParseTolerant(vStr); err == nil {
		return version, true
	}
	return semver.Version{}, false
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

func (container Container) Copy() Container {
	var cp = container
	cp.Commands = append([]string{}, cp.Commands...)
	cp.Env = append([]kubeModels.Env{}, cp.Env...)
	cp.ConfigMaps = append([]kubeModels.ContainerVolume{}, cp.ConfigMaps...)
	cp.VolumeMounts = append([]kubeModels.ContainerVolume{}, cp.VolumeMounts...)
	cp.Ports = append([]kubeModels.ContainerPort{}, cp.Ports...)
	return cp
}

func (container Container) Patch(overlay Container) Container {
	var cp = container.Copy()
	if overlay.Limits.CPU > 0 {
		cp.Limits.CPU = overlay.Limits.CPU
	}
	if overlay.Limits.Memory > 0 {
		cp.Limits.Memory = overlay.Limits.Memory
	}
	cp.Image = str.Vector{overlay.Image, cp.Image}.FirstNonEmpty()
	cp.ConfigMaps = mergeMounts(container.ConfigMountsMap(), overlay.ConfigMountsMap())
	cp.VolumeMounts = mergeMounts(container.VolumeMountsMap(), overlay.VolumeMountsMap())
	cp.PutEnvMap(overlay.GetEnvMap())
	return cp
}

func mergeMounts(lefts, rights map[string]kubeModels.ContainerVolume) []kubeModels.ContainerVolume {
	var mounts = make([]kubeModels.ContainerVolume, 0, (len(lefts)+len(rights))/2)
	for rightPath, right := range rights {
		lefts[rightPath] = right
	}
	for _, patched := range lefts {
		mounts = append(mounts, patched)
	}
	return mounts
}
