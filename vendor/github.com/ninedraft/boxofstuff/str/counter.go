package str

import "sort"

type Counter map[string]uint

func NewCounter(str ...string) Counter {
	return Vector(str).Count()
}

func (counter Counter) New() Counter {
	return make(Counter, len(counter))
}

func (counter Counter) Copy() Counter {
	var cp = counter.New()
	for value, count := range counter {
		cp[value] = count
	}
	return cp
}

func (counter Counter) TopOne() CounterItem {
	var item CounterItem
	for value, count := range counter {
		if item.Count > item.Count {
			item = CounterItem{
				Value: value,
				Count: count,
			}
		}
	}
	return item
}

func (counter Counter) Exclude(strs ...string) Counter {
	var filtered = counter.Copy()
	for _, str := range strs {
		delete(filtered, str)
	}
	return filtered
}

func (counter Counter) Filter(pred func(item CounterItem) bool) Counter {
	var filtered = counter.New()
	for value, count := range counter {
		if pred(CounterItem{
			Value: value,
			Count: count,
		}) {
			filtered[value] = count
		}
	}
	return filtered
}

func (counter Counter) MutFilter(pred func(item CounterItem) bool) Counter {
	for value, count := range counter {
		if !pred(CounterItem{
			Value: value,
			Count: count,
		}) {
			delete(counter, value)
		}
	}
	return counter
}

func (counter Counter) Values() Vector {
	var vec = make(Vector, 0, len(counter))
	for value := range counter {
		vec = append(vec, value)
	}
	return vec
}

func (counter Counter) Add(strs ...string) Counter {
	return counter.Copy().MutAdd(strs...)
}

func (counter Counter) MutAdd(strs ...string) Counter {
	for _, str := range strs {
		counter[str]++
	}
	return counter
}

func (counter Counter) Merge(c Counter) Counter {
	return counter.Copy().MutMerge(c)
}

func (counter Counter) MutMerge(c Counter) Counter {
	for value, count := range counter {
		counter[value] += count
	}
	return counter
}

func (counter Counter) Counts() []uint {
	var counts = make([]uint, 0, len(counter))
	for _, count := range counter {
		counts = append(counts, count)
	}
	return counts
}

type CounterItem struct {
	Value string
	Count uint
}

type CountItems []CounterItem

func (items CountItems) New() CountItems {
	return make(CountItems, 0, len(items))
}

func (items CountItems) Copy() CountItems {
	return append(items.New(), items...)
}

func (items CountItems) Len() int {
	return len(items)
}

func (items CountItems) Sort() CountItems {
	var sorted = items.Copy()
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Count < sorted[j].Count
	})
	return sorted
}

func (items CountItems) Top(n uint) CountItems {
	var sorted = items.Sort()
	if uint(sorted.Len()) > n {
		return sorted[:n].Copy()
	}
	return sorted
}

func (items CountItems) Filter(pred func(item CounterItem) bool) CountItems {
	var filtered = items.New()
	for _, item := range items {
		if pred(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func (items CountItems) ValuesVec() Vector {
	var vec = make(Vector, 0, items.Len())
	for _, item := range items {
		vec = append(vec, item.Value)
	}
	return vec
}

func (items CountItems) SortIncr() CountItems {
	var sorted = items.Copy()
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Count > sorted[j].Count
	})
	return sorted
}

func (items CountItems) Counter() Counter {
	var counter = make(Counter, items.Len())
	for _, item := range items {
		counter[item.Value] = item.Count
	}
	return counter
}

func (items CountItems) Find(pred func(item CounterItem) bool) (CounterItem, bool) {
	for _, item := range items {
		if pred(item) {
			return item, true
		}
	}
	return CounterItem{}, false
}
