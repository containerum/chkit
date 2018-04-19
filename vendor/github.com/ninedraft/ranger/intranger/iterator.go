package intranger

import "strconv"

type intRangerIterator struct {
	intRanger
	value   int
	step    int
	counter int
}

func (iter intRangerIterator) String() string {
	return strconv.Itoa(iter.Value())
}

func (iter *intRangerIterator) Next() bool {
	min, max := iter.Bounds()
	nextVal := iter.counter*iter.step + min
	if nextVal <= max {
		iter.value = nextVal
		iter.counter++
		return true
	}
	return false
}

func (iter intRangerIterator) Value() int {
	return iter.value
}

