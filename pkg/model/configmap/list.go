package configmap

import kubeModels "github.com/containerum/kube-client/pkg/model"

type ConfigMapList []ConfigMap

func ConfigMapListFromKube(kubeList []kubeModels.ConfigMap) ConfigMapList {
	var list = make([]ConfigMap, 0, len(kubeList))
	for _, cm := range kubeList {
		list = append(list, ConfigMapFromKube(cm))
	}
	return list
}

func (list ConfigMapList) ToKube() []kubeModels.ConfigMap {
	var kubeList = make([]kubeModels.ConfigMap, 0, list.Len())
	for _, config := range list {
		kubeList = append(kubeList, config.ToKube())
	}
	return kubeList
}

func (list ConfigMapList) Len() int {
	return len(list)
}

func (list ConfigMapList) Copy() ConfigMapList {
	var cp = append(make(ConfigMapList, 0, list.Len()), list...)
	for i, cm := range cp {
		cp[i] = cm.Copy()
	}
	return cp
}

func (list ConfigMapList) Append(configs ...ConfigMap) ConfigMapList {
	for i := range configs {
		configs[i] = configs[i].Copy()
	}
	return append(list.Copy(), configs...)
}

func (list ConfigMapList) Delete(i int) ConfigMapList {
	list = list.Copy()
	list = append(list[:i], list[i+1:]...)
	return list
}

func (list ConfigMapList) Names() []string {
	var names = make([]string, 0, list.Len())
	for _, config := range list {
		names = append(names, config.Name)
	}
	return names
}

func (list ConfigMapList) Filter(pred func(ConfigMap) bool) ConfigMapList {
	var filtered = make(ConfigMapList, 0, list.Len())
	for _, config := range list {
		if pred(config.Copy()) {
			filtered = append(filtered, config.Copy())
		}
	}
	return filtered
}
