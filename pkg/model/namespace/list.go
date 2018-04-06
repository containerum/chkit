package namespace

import (
	kubeModels "git.containerum.net/ch/kube-client/pkg/model"
)

type NamespaceList []Namespace

func NamespaceListFromKube(kubeList []kubeModels.Namespace) NamespaceList {
	var list NamespaceList = make([]Namespace, 0, len(kubeList))
	for _, namespace := range kubeList {
		list = append(list, NamespaceFromKube(namespace))
	}
	return list
}
