package access

import (
	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type Access kubeModels.UserAccess

func AccessFromKube(kubeAccess kubeModels.UserAccess) Access {
	return Access(kubeAccess)
}

func (access Access) ToKube() kubeModels.UserAccess {
	return kubeModels.UserAccess(access)
}

func (Access) TableHeaders() []string {
	return []string{"Username", "Level"}
}

func (access Access) TableRows() [][]string {
	return [][]string{{
		access.Username,
		access.AccessLevel.String(),
	}}
}

func (access Access) String() string {
	return kubeModels.UserAccess(access).String()
}
