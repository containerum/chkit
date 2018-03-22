package pod

import (
	"git.containerum.net/ch/kube-client/pkg/model"
)

type PodList []Pod

func PodListFromKube(kubeList []model.Pod) PodList {
	var podList PodList = make([]Pod, 0, len(kubeList))
	for _, kubePod := range kubeList {
		podList = append(podList, PodFromKube(kubePod))
	}
	return podList
}
