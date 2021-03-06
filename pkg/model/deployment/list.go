package deployment

import (
	"sort"

	"github.com/blang/semver"
	"github.com/containerum/kube-client/pkg/model"
)

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

func (list DeploymentList) Versions() []semver.Version {
	var versions = make([]semver.Version, 0, list.Len())
	for _, depl := range list {
		versions = append(versions, depl.Version)
	}
	return versions
}

func (list DeploymentList) SortByLess(less func(a, b Deployment) bool) DeploymentList {
	var sorted = list.Copy()
	sort.Slice(sorted, func(i, j int) bool {
		return less(sorted[i].Copy(), sorted[j].Copy())
	})
	return sorted
}

func (list DeploymentList) Inactive() DeploymentList {
	return list.Filter(func(depl Deployment) bool {
		return !depl.Active
	})
}

func (list DeploymentList) Active() DeploymentList {
	return list.Filter(func(depl Deployment) bool {
		return depl.Active
	})
}
