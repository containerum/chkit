package configmap

import (
	"time"

	"github.com/containerum/chkit/pkg/model"
	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type ConfigMap kubeModels.ConfigMap

func ConfigMapFromKube(kubeConfigMap kubeModels.ConfigMap) ConfigMap {
	return ConfigMap(kubeConfigMap).Copy()
}

func (config ConfigMap) ToKube() kubeModels.ConfigMap {
	return kubeModels.ConfigMap(config.Copy())
}

func (config ConfigMap) Copy() ConfigMap {
	var cm = config
	cm.Data = make(kubeModels.ConfigMapData, len(config.Data))
	for k, v := range config.Data {
		cm.Data[k] = v
	}
	return config
}

func (config ConfigMap) Set(key string, value string) ConfigMap {
	config = config.Copy()
	config.Data[key] = value
	return config
}

func (config ConfigMap) Add(data map[string]string) ConfigMap {
	config = config.Copy()
	for k, v := range data {
		config.Data[k] = v
	}
	return config
}

func (config ConfigMap) AddItems(items ...Item) ConfigMap {
	config = config.Copy()
	for _, item := range items {
		config.Data[item.Key] = item.Value
	}
	return config
}

func (config ConfigMap) Items() Items {
	var items = make(Items, 0, len(config.Data))
	for k, v := range config.Data {
		items = append(items, Item{
			Key:   k,
			Value: v,
		})
	}
	return items.Sorted()
}

//  Get -- if defaultValues passed, then first return
func (config ConfigMap) Get(key string, defaultValues ...string) (string, bool) {
	value, ok := config.Data[key]
	if !ok {
		for _, defaultValue := range defaultValues {
			return defaultValue, ok
		}
	}
	return value, ok
}

func (config ConfigMap) Delete(key string) ConfigMap {
	config = config.Copy()
	delete(config.Data, key)
	return config
}

func (config ConfigMap) SetName(name string) ConfigMap {
	config.Name = name
	return config
}

func (config ConfigMap) Age() string {
	if timestamp, err := time.Parse(model.TimestampFormat, config.CreatedAt); err == nil {
		return model.Age(timestamp)
	}
	return "undefined"
}

func (config ConfigMap) New() ConfigMap {
	return ConfigMap{
		Data: make(kubeModels.ConfigMapData, len(config.Data)),
	}
}
