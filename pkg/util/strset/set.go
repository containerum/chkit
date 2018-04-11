package strset

type Set map[string]struct{}

func NewSet(vals []string) Set {
	var s Set = make(map[string]struct{}, len(vals))
	for _, str := range vals {
		s[str] = struct{}{}
	}
	return s
}
func (s Set) Have(v string) bool {
	_, ok := s[v]
	return ok
}

func (s Set) Slice() []string {
	slice := make([]string, 0, len(s))
	for item := range s {
		slice = append(slice, item)
	}
	return slice
}
