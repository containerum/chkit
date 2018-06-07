package deployment

import "github.com/containerum/kube-client/pkg/model"

type DeploymentList []Deployment

func DeploymentListFromKube(kubeList model.DeploymentsList) DeploymentList {
	list := make([]Deployment, 0, len(kubeList.Deployments))
	for _, kubeDeployment := range kubeList.Deployments {
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
