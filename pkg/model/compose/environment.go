package compose

type Environment []string

func (environment Environment) Copy() Environment {
	var cp = make(Environment, 0, len(environment))
	return append(cp, environment...)
}
