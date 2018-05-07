package solution

import (
	"strings"

	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type SolutionList kubeModels.AvailableSolutionsList

func SolutionListFromKube(kubeList kubeModels.AvailableSolutionsList) SolutionList {
	return SolutionList(kubeList)
}

func (list SolutionList) Len() int {
	return kubeModels.AvailableSolutionsList(list).Len()
}

func (list SolutionList) Filter(pred func(Solution) bool) SolutionList {
	return SolutionList(kubeModels.AvailableSolutionsList(list).Filter(func(solution kubeModels.AvailableSolution) bool {
		return pred(Solution(solution))
	}))
}

func (list SolutionList) SearchByName(name string) SolutionList {
	name = strings.ToLower(name)
	return list.Filter(func(solution Solution) bool {
		return strings.Contains(strings.ToLower(solution.Name), name)
	})
}

func (list SolutionList) String() string {
	var strs = make([]string, 0, list.Len())
	for _, sol := range list.Solutions {
		strs = append(strs, SolutionFromKube(sol).String())
	}
	return strings.Join(strs, "\n")
}
