package service

import (
	"git.containerum.net/ch/kube-client/pkg/model"
)

type ServiceList []Service

func ServiceListFromKube(kubeList []model.Service) ServiceList {
	var list ServiceList = make([]Service, 0, len(kubeList))
	for _, kubeService := range kubeList {
		list = append(list, ServiceFromKube(kubeService))
	}
	return list
}

func (list ServiceList) Names() []string {
	names := make([]string, 0, len(list))
	for _, serv := range list {
		names = append(names, serv.Name)
	}
	return names
}

func (list ServiceList) GetByName(name string) (Service, bool) {
	for _, serv := range list {
		if serv.Name == name {
			return serv, true
		}
	}
	return Service{}, false
}
