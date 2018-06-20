package solution

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/util/text"
	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type SolutionTemplate kubeModels.AvailableSolution

func SolutionTemplateFromKube(kubeSolution kubeModels.AvailableSolution) SolutionTemplate {
	return SolutionTemplate(kubeSolution.Copy())
}

func (solution SolutionTemplate) ToKube() kubeModels.AvailableSolution {
	return kubeModels.AvailableSolution(solution).Copy()
}

func (solution SolutionTemplate) Copy() SolutionTemplate {
	return SolutionTemplateFromKube(kubeModels.AvailableSolution(solution))
}

func (solution SolutionTemplate) String() string {
	return fmt.Sprintf(`%s <%s> [%s]`,
		solution.Name, solution.URL, solution.ImagePreview())
}

func (solution SolutionTemplate) ImagePreview() string {
	view := strings.Join(solution.Images, ", ")
	width := text.Width(view)
	const max = 32
	if width > max {
		view = string([]rune(view)[:max]) + "..."
	}
	return view
}

func (solution SolutionTemplate) Describe() string {
	imgs := "Images:"
	return fmt.Sprintf("Name: %s"+
		"URL : %s"+
		imgs+"\n%s",
		solution.Name,
		solution.URL,
		text.Indent(strings.Join(solution.Images, "\n"), uint(len(imgs))))
}

type SolutionTemplates []SolutionTemplate

func (list SolutionTemplates) Len() int {
	return len(list)
}

func (list SolutionTemplates) New() SolutionTemplates {
	return make(SolutionTemplates, 0, list.Len())
}

func (list SolutionTemplates) Copy() SolutionTemplates {
	var cp = list.New()
	for _, templ := range list {
		cp = append(cp, templ.Copy())
	}
	return cp
}

func (list SolutionTemplates) Names() []string {
	var names = make([]string, list.Len())
	for _, templ := range list {
		names = append(names, templ.Name)
	}
	return names
}

func (list SolutionTemplates) GetByName(name string) (SolutionTemplate, bool) {
	for _, templ := range list {
		if templ.Name == name {
			return templ.Copy(), true
		}
	}
	return SolutionTemplate{}, false
}
