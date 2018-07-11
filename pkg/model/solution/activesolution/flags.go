package activesolution

import (
	"errors"

	"github.com/containerum/chkit/pkg/model/solution"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/containerum/chkit/pkg/util/pairs"
)

type Flags struct {
	Force    bool   `flag:"force f" desc:"suppress confirmation, optional"`
	Name     string `desc:"solution name, optional"`
	Template string `desc:"solution template, optional"`
	Env      string `desc:"solution environment variables, optional"`
	Branch   string `desc:"solution git repo branch, optional"`
}

func (flags Flags) Solution(nsID string, args []string) (solution.Solution, error) {
	var sol = solution.Solution{
		Name:      flags.Name,
		Namespace: nsID,
		Branch:    flags.Branch,
		Template:  flags.Template,
	}
	if len(args) == 1 {
		sol.Template = args[0]
	} else if flags.Force {
		//TODO
		return sol, errors.New("not enough arguments")
	}

	if sol.Branch == "" {
		sol.Branch = "master"
	}

	if flags.Name == "" {
		sol.Name = namegen.ColoredPhysics()
		if sol.Template != "" {
			sol.Name += "-" + sol.Template
		}
	}

	if len(sol.Env) != 0 {
		env, err := pairs.ParseMap(flags.Env, ":")
		if err != nil {
			return sol, err
		}
		sol.Env = env
	}
	return sol, nil
}
