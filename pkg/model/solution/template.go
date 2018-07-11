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

func (solution SolutionTemplate) Resources() string {
	return fmt.Sprintf("CPU    : %v mCPU\n"+
		"MEMORY : %v Mb",
		solution.Limits.CPU, solution.Limits.RAM)
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
