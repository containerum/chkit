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

func (list DeploymentList) Len() int {
	return len(list)
}

func (list DeploymentList) New() DeploymentList {
	return make(DeploymentList, 0, len(list))
}

func (list DeploymentList) Copy() DeploymentList {
	var cp = list.New()
	for _, depl := range list {
		cp = append(cp, depl.Copy())
	}
	return cp
}

func (list DeploymentList) Filter(pred func(depl Deployment) bool) DeploymentList {
	var filtered = list.New()
	for _, depl := range list {
		if pred(depl.Copy()) {
			filtered = append(filtered, depl.Copy())
		}
	}
	return filtered
}

func (list DeploymentList) GetByName(name string) (Deployment, bool) {
	for _, depl := range list {
		if depl.Name == name {
			return depl.Copy(), true
		}
	}
	return Deployment{}, false
}
