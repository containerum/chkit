package deployment

import "git.containerum.net/ch/kube-client/pkg/model"

type DeploymentList []Deployment

func DeploymentListFromKube(kubeDeployment []model.Deployment) DeploymentList {
	list := make([]Deployment, 0, len(kubeDeployment))
	return list
}
