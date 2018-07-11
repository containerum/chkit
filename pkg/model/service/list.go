package service

import "github.com/containerum/kube-client/pkg/model"

type ServiceList []Service

func ServiceListFromKube(kubeList model.ServicesList) ServiceList {
	var list ServiceList = make([]Service, 0, len(kubeList.Services))
	for _, kubeService := range kubeList.Services {
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

func (list ServiceList) AvailableForIngress() ServiceList {
	var sortedList ServiceList = make([]Service, 0)
	for _, svc := range list {
		for _, port := range svc.Ports {
			if port.Protocol == "TCP" {
				sortedList = append(sortedList, svc)
				break
			}
		}
	}
	return sortedList
}
