package ingress

import kubeModels "github.com/containerum/kube-client/pkg/model"

type IngressList []Ingress

func IngressListFromKube(kubeList []kubeModels.Ingress) IngressList {
	var list IngressList = make([]Ingress, 0, len(kubeList))
	for _, kubeIngress := range kubeList {
		list = append(list, IngressFromKube(kubeIngress))
	}
	return list
}

func (list IngressList) ToKube() []kubeModels.Ingress {
	var kubeList = make([]kubeModels.Ingress, 0, len(list))
	for _, ingr := range list {
		kubeList = append(kubeList, ingr.ToKube())
	}
	return kubeList
}

func (list IngressList) Copy() IngressList {
	var cp IngressList = make([]Ingress, 0, len(list))
	for _, ingr := range list {
		cp = append(cp, ingr.Copy())
	}
	return cp
}

func (list IngressList) Append(ing ...Ingress) IngressList {
	return append(list.Copy(), ing...)
}

func (list IngressList) Delete(i int) IngressList {
	cp := list.Copy()
	return append(cp[:i], cp[i+1:]...)
}
