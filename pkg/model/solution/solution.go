package solution

import kubeModels "github.com/containerum/kube-client/pkg/model"

type UserSolution kubeModels.UserSolution

func UserolutionFromKube(kubeSolution kubeModels.UserSolution) UserSolution {
	return UserSolution(kubeSolution)
}

func (solution UserSolution) ToKube() kubeModels.UserSolution {
	return kubeModels.UserSolution(solution).Copy()
}
