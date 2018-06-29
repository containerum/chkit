package activesolution

import (
	"io/ioutil"
	"path"

	"bytes"
	"encoding/json"
	"os"

	"errors"

	"github.com/containerum/chkit/pkg/model/solution"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/containerum/chkit/pkg/util/pairs"
	"gopkg.in/yaml.v2"
)

type Flags struct {
	Force     bool   `flag:"force f" desc:"suppress confirmation, optional"`
	File      string `desc:"file with solution data, .yaml or .json, stdin if '-', optional"`
	Name      string `desc:"solution name, optional"`
	Namespace string `desc:"solution namespace, optional"`
	Template  string `desc:"solution template, optional"`
	Env       string `desc:"solution environment variables, optional"`
	Branch    string `desc:"solution git repo branch, optional"`
}

func (flags Flags) Solution(nsID string, args []string) (solution.Solution, error) {
	var sol = solution.Solution{
		Name:      flags.Name,
		Namespace: nsID,
		Branch:    flags.Branch,
		Template:  flags.Template,
	}
	if flags.File != "" {
		var err error
		sol, err = flags.solutionFromFile()
		if err != nil {
			ferr.Println(err)
			return sol, err
		}
	} else if len(args) == 1 {
		sol.Template = args[0]
	} else if flags.Force {
		//TODO
		return sol, errors.New("not enough arguments")
	}

	if sol.Branch == "" {
		sol.Branch = "master"
	}

	if flags.Namespace != "" {
		sol.Namespace = flags.Namespace
	}

	if flags.Name == "" {
		if sol.Template != "" {
			sol.Name = namegen.Color() + "-" + sol.Template
		} else {
			sol.Name = namegen.ColoredPhysics()
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

func (flags Flags) solutionFromFile() (solution.Solution, error) {
	var sol solution.Solution
	data, err := func() ([]byte, error) {
		if flags.File == "-" {
			buf := &bytes.Buffer{}
			_, err := buf.ReadFrom(os.Stdin)
			if err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		}
		data, err := ioutil.ReadFile(flags.File)
		if err != nil {
			return data, err
		}
		return data, nil
	}()
	if err != nil {
		return sol, err
	}
	if path.Ext(flags.File) == "yaml" {
		if err := yaml.Unmarshal(data, &sol); err != nil {
			return sol, err
		} else {
			if err := json.Unmarshal(data, &sol); err != nil {
				return sol, err
			}

		}
	}
	return sol, nil
}
