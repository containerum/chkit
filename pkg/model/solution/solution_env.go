package solution

import (
	"fmt"

	"sort"

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

type Env struct {
	Name  string
	Value string
}

func (env Env) String() string {
	return env.Name + ":" + env.Value
}

type Envs []Env

func EnvsFromMap(m map[string]string) Envs {
	var envs = make(Envs, 0, len(m))
	for name, value := range m {
		envs = append(envs, Env{
			Name:  name,
			Value: value,
		})
	}
	return envs
}

func (envs Envs) New() Envs {
	return make(Envs, 0, len(envs))
}

func (envs Envs) Copy() Envs {
	return append(envs.New(), envs...)
}

func (envs Envs) Sorted() Envs {
	var sorted = envs.Copy()
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})
	return sorted
}

func (envs Envs) Map() map[string]string {
	var m = make(map[string]string, len(envs))
	for _, env := range envs {
		m[env.Name] = env.Value
	}
	return m
}
