package solution

import (
	"fmt"

	"github.com/containerum/chkit/pkg/util/text"
	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type Solution kubeModels.UserSolution

func SolutionFromKube(kubeSolution kubeModels.UserSolution) Solution {
	return Solution(kubeSolution)
}

func (solution Solution) ToKube() kubeModels.UserSolution {
	return kubeModels.UserSolution(solution).Copy()
}

//TODO
func (solution Solution) String() string {
	return fmt.Sprintf(`%s <%s> [%s]`,
		solution.Name, solution.Template, solution.Namespace)
}

func (solution Solution) Describe() string {
	envst := "ENV:"

	env := kubeModels.SolutionEnv{Env: solution.Env}

	return fmt.Sprintf("Name: %s"+
		"URL : %s"+
		envst+"\n%s",
		solution.Name,
		solution.URL,
		text.Indent(SolutionEnvFromKube(env).String(), uint(len(env.Env))))
}
