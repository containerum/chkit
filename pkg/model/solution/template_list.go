package solution

import (
	"strings"

	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type TemplatesList kubeModels.AvailableSolutionsList

func TemplatesListFromKube(kubeList kubeModels.AvailableSolutionsList) TemplatesList {
	return TemplatesList(kubeList)
}

func (list TemplatesList) Len() int {
	return kubeModels.AvailableSolutionsList(list).Len()
}

func (list TemplatesList) Filter(pred func(template SolutionTemplate) bool) TemplatesList {
	return TemplatesList(kubeModels.AvailableSolutionsList(list).Filter(func(solution kubeModels.AvailableSolution) bool {
		return pred(SolutionTemplate(solution))
	}))
}

func (list TemplatesList) SearchByName(name string) TemplatesList {
	name = strings.ToLower(name)
	return list.Filter(func(solution SolutionTemplate) bool {
		return strings.Contains(strings.ToLower(solution.Name), name)
	})
}

func (list TemplatesList) String() string {
	var strs = make([]string, 0, list.Len())
	for _, sol := range list.Solutions {
		strs = append(strs, SolutionTemplateFromKube(sol).String())
	}
	return strings.Join(strs, "\n")
}

func (list TemplatesList) Names() []string {
	var names = make([]string, 0, list.Len())
	for _, sol := range list.Solutions {
		names = append(names, sol.Name)
	}
	return names
}
