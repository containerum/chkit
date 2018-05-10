package configmap

import "fmt"

type Item struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func NewItem(key string, value interface{}) Item {
	return Item{
		Key:   key,
		Value: value,
	}
}

func (item Item) Data() (key string, value interface{}) {
	return item.Key, item.Value
}

func (item Item) String() string {
	return fmt.Sprintf("%s : %v", item.Key, item.Value)
}
