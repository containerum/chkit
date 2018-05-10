package configmap

import (
	"fmt"
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
	cm.Data = make(map[string]interface{}, len(config.Data))
	for k, v := range config.Data {
		cm.Data[k] = v
	}
	return config
}

func (config ConfigMap) Set(key string, value interface{}) ConfigMap {
	config = config.Copy()
	config.Data[key] = value
	return config
}

func (config ConfigMap) Add(data map[string]interface{}) ConfigMap {
	config = config.Copy()
	for k, v := range data {
		config.Data[k] = v
	}
	return config
}

func (config ConfigMap) Get(key string, defaultValues ...interface{}) (interface{}, bool) {
	value, ok := config.Data[key]
	if !ok {
		for _, defaultValue := range defaultValues {
			if defaultValue != nil {
				return defaultValue, ok
			}
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
	if config.CreatedAt == nil {
		return "undefined"
	}
	timestamp, err := time.Parse(time.RFC3339, *config.CreatedAt)
	if err != nil {
		return fmt.Sprintf("invlalid timestamp %q", *config.CreatedAt)
	}
	return model.Age(timestamp)
}
