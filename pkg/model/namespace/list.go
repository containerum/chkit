package namespace

import (
	"fmt"
	"strings"

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

func (list NamespaceList) GetDefault(i int, defaultNs Namespace) (Namespace, bool) {
	if i >= 0 && i < list.Len() {
		return list.Get(i), true
	}
	return defaultNs, false
}

func (list NamespaceList) Get(i int) Namespace {
	return list[i].Copy()
}

func (list NamespaceList) Head() (Namespace, bool) {
	return list.GetDefault(0, Namespace{})
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

func (list NamespaceList) Filter(pred func(Namespace) bool) NamespaceList {
	var filtered = list.New()
	for _, namespace := range list {
		if pred(namespace.Copy()) {
			filtered = append(filtered, namespace.Copy())
		}
	}
	return filtered
}

// get Namespace by string $LABEL or $OWNER_LOGIN/$LABEL
func (list NamespaceList) GetByUserFriendlyID(label string) (Namespace, bool) {
	var tokens = strings.SplitN(label, "/", 2)
	if len(tokens) == 2 {
		return list.GetByLabelAndOwner(tokens[0], tokens[1])
	}
	return list.Filter(func(namespace Namespace) bool {
		return namespace.Label == label
	}).Head()
}

func (list NamespaceList) GetByLabelAndOwner(owner, label string) (Namespace, bool) {
	return list.Filter(func(namespace Namespace) bool {
		return namespace.OwnerLogin == owner && namespace.Label == label
	}).Head()
}
