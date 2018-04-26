package intranger

import (
	"fmt"
	"strconv"
)

type intRanger struct {
	min int
	max int
}

func IntRanger(min, max int) intRanger {
	if min >= max {
		panic(fmt.Sprintf("[ranger IntRanger] min parameter must be less than max, but min=%d, max=%d", min, max))
	}
	return intRanger{
		min: min,
		max: max,
	}
}

func DefaultIntRanger() intRanger {
	return IntRanger(0, 1)
}

func (ranger intRanger) Containing(i int) bool {
	return ranger.min <= i && i <= ranger.max
}

func (ranger intRanger) String() string {
	return strconv.Itoa(ranger.min) + ".." + strconv.Itoa(ranger.max)
}

func (ranger intRanger) Min() int {
	return ranger.min
}

func (ranger intRanger) Max() int {
	return ranger.max
}

func (ranger intRanger) Bounds() (min int, max int) {
	return ranger.Min(), ranger.Max()
}

func (ranger intRanger) In(r intRanger) bool {
	return ranger.min >= r.min && ranger.max <= r.max
}


func (ranger intRanger) Iter() intRangerIterator {
	return ranger.IterWithStep(1)
}

func (ranger intRanger) IterWithStep(step int) intRangerIterator {
	return intRangerIterator{
		intRanger: ranger,
		value:     ranger.min,
		step:      step,
		counter:   0,
	}
}

func (ranger intRanger) Shift(off int) intRanger {
	return IntRanger(ranger.min+off, ranger.max+off)
}
