package deployment

import "git.containerum.net/ch/kube-client/pkg/model"

type DeploymentList []Deployment

func DeploymentListFromKube(kubeList []model.Deployment) DeploymentList {
	list := make([]Deployment, 0, len(kubeList))
	for _, kubeDeployment := range kubeList {
		list = append(list, DeploymentFromKube(kubeDeployment))
	}
	return list
}

func (list DeploymentList) Names() []string {
	names := make([]string, 0, len(list))
	for _, depl := range list {
		names = append(names, depl.Name)
	}
	return names
}
