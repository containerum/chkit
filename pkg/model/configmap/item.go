package configmap

import (
	"encoding/base64"
	"encoding/json"
	"sort"
	"fmt"
)

type Item struct {
	key   string
	value string
}

func (item Item) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(item.toJSON(), "", "  ")
}

func (item *Item) UnmarshalJSON(b []byte) error {
	var i _jsonItem
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	*item = fromJSON(i)
	return nil
}

type _jsonItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (item Item) toJSON() _jsonItem {
	return _jsonItem{
		Key:   item.key,
		Value: item.value,
	}
}

func fromJSON(jsItem _jsonItem) Item {
	return Item{
		key:   jsItem.Key,
		value: jsItem.Value,
	}
}

func NewItem(key string, value string) Item {
	var encodedValue = base64.StdEncoding.EncodeToString([]byte(value))
	return Item{
		key:   key,
		value: encodedValue,
	}
}

func (item Item) Key() string {
	return item.key
}

func (item Item) Value() string {
	var decodedValue, err = base64.StdEncoding.DecodeString(item.value)
	if err == nil {
		return string(decodedValue)
	}
	return item.value
}

func (item Item) Data() (key string, value string) {
	return item.key, item.Value()
}

func (item Item) String() string {
	var key, value = item.Data()
	return key + ":" + value
}

func (item Item) WithKey(key string) Item {
	return Item{
		key:   key,
		value: item.value,
	}
}

func (item Item) WithValue(value string) Item {
	return NewItem(
		item.key,
		value,
	)
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
