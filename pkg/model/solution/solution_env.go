package solution

import (
	"fmt"

	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type SolutionEnv kubeModels.SolutionEnv

func SolutionEnvFromKube(kubeSolutionEnv kubeModels.SolutionEnv) SolutionEnv {
	return SolutionEnv(kubeSolutionEnv)
}

func (solutionEnv SolutionEnv) String() string {
	var ret string
	for k, v := range solutionEnv.Env {
		ret += fmt.Sprintf(`%s = %s;`, k, v)
	}
	return ret
}
