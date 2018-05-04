package solution

import (
	"strings"

	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type SolutionList kubeModels.AvailableSolutionsList

func SolutionListFromKube(kubeList SolutionList) SolutionList {
	return SolutionList(kubeList)
}

func (list SolutionList) Filter(pred func(Solution) bool) SolutionList {
	return SolutionList(kubeModels.AvailableSolutionsList(list).Filter(func(solution kubeModels.AvailableSolution) bool {
		return pred(Solution(solution))
	}))
}

func (list SolutionList) SearchByName(name string) SolutionList {
	return list.Filter(func(solution Solution) bool {
		return strings.Contains(solution.Name, name)
	})
}
