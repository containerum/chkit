package configmap

import (
	"fmt"
	"sort"
)

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewItem(key string, value string) Item {
	return Item{
		Key:   key,
		Value: value,
	}
}

func (item Item) Data() (key string, value string) {
	return item.Key, item.Value
}

func (item Item) String() string {
	return fmt.Sprintf("%s : %q", item.Key, item.Value)
}

type Items []Item

func (items Items) New() Items {
	return make(Items, 0, len(items))
}

func (items Items) Copy() Items {
	return append(items.New(), items...)
}

func (items Items) Sorted() Items {
	var cp = items.Copy()
	sort.Slice(cp, func(i, j int) bool {
		return cp[i].Key < cp[j].Key
	})
	return cp
}

func (items Items) Map() map[string]string {
	var m = make(map[string]string, len(items))
	for _, item := range items {
		m[item.Key] = item.Value
	}
	return m
}
