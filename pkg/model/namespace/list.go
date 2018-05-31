package namespace

import (
	"fmt"

	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type NamespaceList []Namespace

func NamespaceListFromKube(kubeList kubeModels.NamespacesList) NamespaceList {
	var list NamespaceList = make([]Namespace, 0, len(kubeList.Namespaces))
	for _, namespace := range kubeList.Namespaces {
		list = append(list, NamespaceFromKube(namespace).Copy())
	}
	return list
}

func (list NamespaceList) New() NamespaceList {
	return make(NamespaceList, 0, len(list))
}

func (list NamespaceList) Copy() NamespaceList {
	var cp = list.New()
	for _, namespace := range list {
		cp = append(cp, namespace.Copy())
	}
	return cp
}

func (list NamespaceList) ToKube() kubeModels.NamespacesList {
	var namespaces = make([]kubeModels.Namespace, 0, len(list))
	for _, namespace := range list {
		namespaces = append(namespaces, namespace.ToKube())
	}
	return kubeModels.NamespacesList{
		Namespaces: namespaces,
	}
}

func (list NamespaceList) Len() int {
	return len(list)
}

func (list NamespaceList) Labels() []string {
	var labels = make([]string, 0, list.Len())
	for _, namespace := range list {
		labels = append(labels, namespace.Label)
	}
	return labels
}

func (list NamespaceList) IDs() []string {
	var IDs = make([]string, 0, list.Len())
	for _, namespace := range list {
		IDs = append(IDs, namespace.ID)
	}
	return IDs
}

func (list NamespaceList) LabelsAndIDs() []string {
	var lines = make([]string, 0, list.Len())
	for _, namespace := range list {
		lines = append(lines, fmt.Sprintf("%s %s", namespace.Label, namespace.ID))
	}
	return lines
}
