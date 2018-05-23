package volume

import kubeModels "github.com/containerum/kube-client/pkg/model"

type VolumeList []Volume

func VolumeListFromKube(kubeList []kubeModels.Volume) VolumeList {
	var list VolumeList = make([]Volume, 0, len(kubeList))
	for _, volume := range kubeList {
		list = append(list, VolumeFromKube(volume))
	}
	return list
}
