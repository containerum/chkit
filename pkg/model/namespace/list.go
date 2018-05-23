package namespace

import (
	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type NamespaceList []Namespace

func NamespaceListFromKube(kubeList []kubeModels.Namespace) NamespaceList {
	var list NamespaceList = make([]Namespace, 0, len(kubeList))
	for _, namespace := range kubeList {
		list = append(list, NamespaceFromKube(namespace))
	}
	return list
}
