package str

import "strings"

type Vector []string

func (vector Vector) FirstNonEmpty() string {
	for _, str := range vector {
		if str != "" {
			return str
		}
	}
	return ""
}

func (vector Vector) Len() int {
	return len(vector)
}

func (vector Vector) New() Vector {
	return make(Vector, 0, vector.Len())
}

func (vector Vector) Copy() Vector {
	return append(vector.New(), vector...)
}

func (vector Vector) Filter(pred func(str string) bool) Vector {
	var filtered = vector.New()
	for _, str := range vector {
		if pred(str) {
			filtered = append(filtered, str)
		}
	}
	return filtered
}

func (vector Vector) Count() map[string]int {
	var count = make(map[string]int, vector.Len())
	for _, str := range vector {
		count[str]++
	}
	return count
}

func (vector Vector) Slice() []string {
	return []string(vector.Copy())
}

func (vector Vector) Join(delim string) string {
	return strings.Join(vector, delim)
}

func (vector Vector) WriteToChan(ch chan<- string) {
	for _, str := range vector {
		ch <- str
	}
}

func (vector Vector) Append(strs ...string) Vector {
	return append(vector.Copy(), strs...)
}

func (vector Vector) Map(op func(str string) string) Vector {
	var mapped = vector.New()
	for _, str := range vector {
		mapped = append(mapped, op(str))
	}
	return mapped
}

func (vector Vector) Populate(op func(str string) []string) Vector {
	var populated = vector.New()
	for _, str := range vector {
		populated = append(populated, op(str)...)
	}
	return populated
}

func (vector Vector) Fold(start string, op func(acc, str string) string) string {
	var acc = start
	for _, str := range vector {
		acc = op(acc, str)
	}
	return acc
}

func (vector Vector) Delete(i int) Vector {
	vector = vector.Copy()
	return append(vector[:i], vector[i+1:]...)
}
