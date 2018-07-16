package compose

import (
	"strconv"

	"gopkg.in/yaml.v2"
)

var (
	_ yaml.Unmarshaler = new(Build)
)

type Build struct {
	Context    string
	Dockerfile string
	Args       BuildArgs
	Target     string
	Command    []string
}

func (build Build) Copy() Build {
	var commandCp = make([]string, 0, len(build.Command))
	build.Command = commandCp
	return build
}

func (build *Build) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value interface{}
	if err := unmarshal(&value); err != nil {
		return err
	}
	type yamlBuild Build
	switch value := value.(type) {
	case string:
		*build = Build{Context: value}
	case map[string]interface{}:
		var data = &yamlBuild{}
		if err := unmarshal(&data); err != nil {
			return err
		}
		*build = Build(*data)
	}
	return nil
}

var (
	_ yaml.Unmarshaler = new(BuildArgs)
)

type BuildArgs map[string]string

func (args BuildArgs) New() BuildArgs {
	return make(BuildArgs, len(args))
}

func (args BuildArgs) Copy() BuildArgs {
	var cp = args.New()
	for k, v := range args {
		cp[k] = v
	}
	return cp
}

func (args BuildArgs) Slice() []string {
	var strs = make([]string, 0, len(args))
	for k, v := range args {
		strs = append(strs, k+":"+strconv.Quote(v))
	}
	return strs
}

func (args *BuildArgs) UnmarshalYAML(unmarshal Unmarshaler) error {
	var value interface{}
	if err := unmarshal(&value); err != nil {
		return err
	}
	type yamlArs BuildArgs
	switch value := value.(type) {
	case []string:
		panic("NOT IMPLEMENTED")
	case map[string]string:
		*args = BuildArgs(value)
	}
	return nil
}
