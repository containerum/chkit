package solution

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/util/text"
	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type Solution kubeModels.AvailableSolution

func SolutionFromKube(kubeSolution kubeModels.AvailableSolution) Solution {
	return Solution(kubeSolution.Copy())
}

func (solution Solution) ToKube() kubeModels.AvailableSolution {
	return kubeModels.AvailableSolution(solution).Copy()
}

func (solution Solution) Copy() Solution {
	return SolutionFromKube(kubeModels.AvailableSolution(solution))
}

func (solution Solution) String() string {
	return fmt.Sprintf(`%s <%s> [%s]`,
		solution.Name, solution.URL, solution.ImagePreview())
}

func (solution Solution) ImagePreview() string {
	view := strings.Join(solution.Images, ", ")
	width := text.Width(view)
	if width > 13 {
		view = string([]rune(view)[:13]) + "..."
	}
	return view
}

func (solution Solution) Describe() string {
	imgs := "Images:"
	return fmt.Sprintf("Name: %s"+
		"URL : %s"+
		imgs+"\n%s",
		solution.Name,
		solution.URL,
		text.Indent(strings.Join(solution.Images, "\n"), uint(len(imgs))))
}
