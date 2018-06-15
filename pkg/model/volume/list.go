package volume

import kubeModels "github.com/containerum/kube-client/pkg/model"

type VolumeList []Volume

func VolumeListFromKube(kubeList kubeModels.VolumesList) VolumeList {
	var list = make(VolumeList, 0, len(kubeList.Volumes))
	for _, volume := range kubeList.Volumes {
		list = append(list, VolumeFromKube(volume))
	}
	return list
}

func (list VolumeList) New() VolumeList {
	return make(VolumeList, 0, len(list))
}

func (list VolumeList) Copy() VolumeList {
	var cp = list.New()
	for _, volume := range list {
		cp = append(cp, volume.Copy())
	}
	return cp
}

func (list VolumeList) ToKube() kubeModels.VolumesList {
	var volumes = make([]kubeModels.Volume, 0, len(list))
	for _, volume := range list {
		volumes = append(volumes, volume.ToKube())
	}
	return kubeModels.VolumesList{
		Volumes: volumes,
	}
}

func (list VolumeList) Names() []string {
	var names = make([]string, 0, len(list))
	for _, volume := range list {
		names = append(names, volume.Name)
	}
	return names
}

func (list VolumeList) Filter(pred func(Volume) bool) VolumeList {
	var filtered = list.New()
	for _, volume := range list {
		if pred(volume.Copy()) {
			filtered = append(filtered, volume.Copy())
		}
	}
	return filtered
}
