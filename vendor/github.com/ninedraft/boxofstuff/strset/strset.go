package strset

import "encoding/json"

type Set map[string]struct{}

func (s Set) Copy() Set {
	var cp = make(Set, len(s))
	for k := range s {
		cp[k] = struct{}{}
	}
	return cp
}

func (s Set) New() Set {
	return make(Set, len(s))
}

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

func (s Set) Chan() <-chan string {
	var ch = make(chan string)
	go func() {
		defer close(ch)
		for elem := range s {
			ch <- elem
		}
	}()
	return ch
}

func (s Set) Map(op func(string) string) Set {
	var cp = s.New()
	for elem := range s {
		cp[op(elem)] = struct{}{}
	}
	return cp
}

func (s Set) Delete(elem string) Set {
	s = s.Copy()
	delete(s, elem)
	return s
}

func (s Set) Filter(pred func(string) bool) Set {
	var cp = s.New()
	for elem := range s {
		if pred(elem) {
			cp[elem] = struct{}{}
		}
	}
	return cp
}

func (s Set) Intersect(x Set) Set {
	var result = s.New()
	for elem := range s {
		if x.Have(elem) {
			result[elem] = struct{}{}
		}
	}
	return result
}

func (s Set) Add(x Set) Set {
	s = s.Copy()
	for elem := range x {
		s[elem] = struct{}{}
	}
	return s
}

func (s Set) Sub(x Set) Set {
	s = s.Copy()
	for elem := range x {
		delete(s, elem)
	}
	return s
}

func (s Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Slice())
}

func (s *Set) UnmarshalJSON(p []byte) error {
	var b []string
	if err := json.Unmarshal(p, &b); err != nil {
		return err
	}
	*s = NewSet(b)
	return nil
}
