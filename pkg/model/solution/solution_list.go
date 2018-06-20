package solution

import (
	"strings"

	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type SolutionsList kubeModels.UserSolutionsList

func SolutionsListFromKube(kubeList kubeModels.UserSolutionsList) SolutionsList {
	return SolutionsList(kubeList)
}

func (list SolutionsList) Len() int {
	return kubeModels.UserSolutionsList(list).Len()
}

func (list SolutionsList) Filter(pred func(Solution) bool) SolutionsList {
	return SolutionsList(kubeModels.UserSolutionsList(list).Filter(func(solution kubeModels.UserSolution) bool {
		return pred(Solution(solution))
	}))
}

func (list SolutionsList) SearchByName(name string) SolutionsList {
	name = strings.ToLower(name)
	return list.Filter(func(solution Solution) bool {
		return strings.Contains(strings.ToLower(solution.Name), name)
	})
}

func (list SolutionsList) String() string {
	var strs = make([]string, 0, list.Len())
	for _, sol := range list.Solutions {
		strs = append(strs, SolutionFromKube(sol).String())
	}
	return strings.Join(strs, "\n")
}

func (list SolutionsList) Names() []string {
	var names = make([]string, 0, list.Len())
	for _, sol := range list.Solutions {
		names = append(names, sol.Name)
	}
	return names
}
