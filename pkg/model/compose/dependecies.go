package compose

type Dependecies []string

func (dependecies Dependecies) Copy() Dependecies {
	var cp = make(Dependecies, 0, len(dependecies))
	cp = append(cp, dependecies...)
	return cp
}
