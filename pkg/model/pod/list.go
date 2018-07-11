package pod

import (
	"strings"

	"github.com/containerum/kube-client/pkg/model"
	"github.com/ninedraft/boxofstuff/str"
	"github.com/ninedraft/boxofstuff/strset"
)

type PodList []Pod

func PodListFromKube(kubeList model.PodsList) PodList {
	var podList PodList = make([]Pod, 0, len(kubeList.Pods))
	for _, kubePod := range kubeList.Pods {
		podList = append(podList, PodFromKube(kubePod))
	}
	return podList
}

func (list PodList) Len() int {
	return len(list)
}

func (list PodList) New() PodList {
	return make(PodList, 0, len(list))
}

func (list PodList) Copy() PodList {
	var cp = list.New()
	for _, po := range list {
		cp = append(cp, po.Copy())
	}
	return cp
}

func (list PodList) Filter(pred func(po Pod) bool) PodList {
	var filtered = list.New()
	for _, po := range list {
		if pred(po.Copy()) {
			filtered = append(filtered, po.Copy())
		}
	}
	return filtered
}

func (list PodList) FilterByStatus(status ...string) PodList {
	var statuses = str.Vector(status).Map(strings.ToLower)
	return list.Filter(func(po Pod) bool {
		return statuses.Contains(strings.ToLower(po.Status.Phase))
	})
}

func (list PodList) FilterByNames(names ...string) PodList {
	var nameSet = strset.NewSet(names)
	return list.Filter(func(po Pod) bool {
		return nameSet.Have(po.Name)
	})
}
