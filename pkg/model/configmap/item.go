package configmap

import "fmt"

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
