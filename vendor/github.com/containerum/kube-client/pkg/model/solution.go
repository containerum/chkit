package model

// AvailableSolutionsList -- list of available solutions
//
// swagger:model
type AvailableSolutionsList struct {
	Solutions []AvailableSolution `json:"solutions"`
}

func (list AvailableSolutionsList) Len() int {
	return len(list.Solutions)
}

func (list AvailableSolutionsList) Copy() AvailableSolutionsList {
	var solutions = make([]AvailableSolution, 0, list.Len())
	for _, sol := range solutions {
		solutions = append(solutions, sol.Copy())
	}
	return AvailableSolutionsList{
		Solutions: solutions,
	}
}

func (list AvailableSolutionsList) Get(i int) AvailableSolution {
	return list.Solutions[i]
}

func (list AvailableSolutionsList) Filter(pred func(AvailableSolution) bool) AvailableSolutionsList {
	solutions := make([]AvailableSolution, 0, list.Len())
	for _, sol := range list.Solutions {
		if pred(sol.Copy()) {
			solutions = append(solutions, sol.Copy())
		}
	}
	return AvailableSolutionsList{
		Solutions: solutions,
	}
}

// AvailableSolution -- solution which user can run
//
// swagger:model
type AvailableSolution struct {
	ID     string          `json:"id,omitempty"`
	Name   string          `json:"name"`
	Limits *SolutionLimits `json:"limits"`
	Images []string        `json:"images"`
	URL    string          `json:"url"`
	Active bool            `json:"active,omitempty"`
}

func (solution AvailableSolution) Copy() AvailableSolution {
	return AvailableSolution{
		Name: solution.Name,
		Limits: func() *SolutionLimits {
			if solution.Limits == nil {
				return nil
			}
			cp := *solution.Limits
			return &cp
		}(),
		Images: append([]string{}, solution.Images...),
		URL:    solution.URL,
		Active: solution.Active,
	}
}

// SolutionLimits -- solution resources limits
//
// swagger:model
type SolutionLimits struct {
	CPU string `json:"cpu"`
	RAM string `json:"ram"`
}

// SolutionEnv -- solution environment variables
//
// swagger:model
type SolutionEnv struct {
	Env map[string]string `json:"env"`
}

// SolutionResources -- list of solution resources
//
// swagger:model
type SolutionResources struct {
	Resources map[string]int `json:"resources"`
}

func (res SolutionResources) Copy() SolutionResources {
	r := make(map[string]int, len(res.Resources))
	for k, v := range res.Resources {
		r[k] = v
	}
	return SolutionResources{
		Resources: r,
	}
}

type ConfigFile struct {
	Name string `json:"config_file"`
	Type string `json:"type"`
}

// UserSolutionsList -- list of running solution
//
// swagger:model
type UserSolutionsList struct {
	Solutions []UserSolution `json:"solutions"`
}

func (list UserSolutionsList) Copy() UserSolutionsList {
	var solutions = make([]UserSolution, 0, list.Len())
	for _, s := range solutions {
		solutions = append(solutions, s.Copy())
	}
	return UserSolutionsList{
		Solutions: solutions,
	}
}

func (list UserSolutionsList) Len() int {
	return len(list.Solutions)
}

func (list UserSolutionsList) Get(i int) UserSolution {
	return list.Solutions[i]
}

// UserSolution -- running solution
//
// swagger:model
type UserSolution struct {
	ID     string            `json:"id,omitempty"`
	Branch string            `json:"branch"`
	Env    map[string]string `json:"env"`
	// required: true
	Template string `json:"template"`
	// required: true
	Name string `json:"name"`
	// required: true
	Namespace string `json:"namespace"`
}

func (solution UserSolution) Copy() UserSolution {
	env := make(map[string]string, len(solution.Env))
	for k, v := range solution.Env {
		env[k] = v
	}
	return UserSolution{
		Branch:    solution.Branch,
		Env:       env,
		Template:  solution.Template,
		Name:      solution.Name,
		Namespace: solution.Namespace,
	}
}

// RunSolutionResponse -- response to run solution request
//
// swagger:model
type RunSolutionResponse struct {
	Created    int      `json:"created"`
	NotCreated int      `json:"not_created"`
	Errors     []string `json:"errors,omitempty"`
}
