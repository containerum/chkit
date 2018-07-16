package compose

import "strconv"

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
