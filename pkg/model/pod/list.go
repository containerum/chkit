package pod

import (
	"github.com/containerum/kube-client/pkg/model"
)

type PodList []Pod

func PodListFromKube(kubeList model.PodsList) PodList {
	var podList PodList = make([]Pod, 0, len(kubeList.Pods))
	for _, kubePod := range kubeList.Pods {
		podList = append(podList, PodFromKube(kubePod))
	}
	return podList
}
