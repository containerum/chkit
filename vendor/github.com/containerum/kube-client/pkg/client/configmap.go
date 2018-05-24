package client

import (
	"github.com/containerum/kube-client/pkg/model"
	"github.com/containerum/kube-client/pkg/rest"
)

const (
	configMapsPath = "/namespaces/{namespace}/configmaps"
	configMapPath  = "/namespaces/{namespace}/configmaps/{configmap}"
)

// CreateConfigMap -- creates a ConfigMap in provided namespace.
func (client *Client) CreateConfigMap(namespace, name string, data model.ConfigMapData) error {
	return client.RestAPI.Post(rest.Rq{
		URL: rest.URL{
			Path: configMapsPath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
		Body: model.ConfigMap{
			Name: name,
			Data: data,
		},
	})
}

// GetConfigMap -- retrieves ConfigMap by name from provided namespace.
func (client *Client) GetConfigMap(namespace, name string) (ret model.ConfigMap, err error) {
	err = client.RestAPI.Get(rest.Rq{
		Result: &ret,
		URL: rest.URL{
			Path: configMapPath,
			Params: rest.P{
				"namespace": namespace,
				"configmap": name,
			},
		},
	})
	return
}

// GetConfigMapList -- returns all ConfigMap`s in namespace.
func (client *Client) GetConfigMapList(namespace string) (ret []model.ConfigMap, err error) {
	jsonAdaptor := struct {
		ConfigMaps *[]model.ConfigMap `json:"configmaps"`
	}{&ret}
	err = client.RestAPI.Get(rest.Rq{
		Result: &jsonAdaptor,
		URL: rest.URL{
			Path: configMapsPath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
	return
}

// UpdateConfigMap -- rewrites ConfigMap by name in provided namespace.
func (client *Client) UpdateConfigMap(namespace, name string, data model.ConfigMapData) error {
	return client.RestAPI.Put(rest.Rq{
		URL: rest.URL{
			Path: configMapsPath,
			Params: rest.P{
				"namespace": namespace,
				"configmap": name,
			},
		},
		Body: model.ConfigMap{
			Data: data,
		},
	})
}

// DeleteConfigMap -- deletes ConfigMap by name in provided namespace
func (client *Client) DeleteConfigMap(namespace, name string) error {
	return client.RestAPI.Delete(rest.Rq{
		URL: rest.URL{
			Path: configMapPath,
			Params: rest.P{
				"namespace": namespace,
				"configmap": name,
			},
		},
	})
}
