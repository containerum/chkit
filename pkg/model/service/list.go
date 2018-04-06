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
