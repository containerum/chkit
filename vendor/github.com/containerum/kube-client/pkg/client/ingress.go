package client

import (
	"github.com/containerum/kube-client/pkg/model"
	"github.com/containerum/kube-client/pkg/rest"
)

const (
	kubeAPIIngressRootPath  = "/namespaces/{namespace}/ingresses"
	kubeAPIIngressPath      = kubeAPIIngressRootPath + "/{domain}"
	resourceIngressRootPath = "/namespace/{namespace}/ingress"
	resourceIngressPath     = resourceIngressRootPath + "/{domain}"
)

// AddIngress -- adds ingress to provided namespace
func (client *Client) AddIngress(namespace string, ingress model.Ingress) error {
	return client.RestAPI.Post(rest.Rq{
		Body: ingress,
		URL: rest.URL{
			Path: resourceIngressRootPath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
}

// GetIngressList -- returns list of ingresses.
func (client *Client) GetIngressList(namespace string) ([]model.Ingress, error) {
	var ingressList []model.Ingress
	jsonAdaptor := struct {
		Ingresses *[]model.Ingress `json:"ingresses"`
	}{&ingressList}
	err := client.RestAPI.Get(rest.Rq{
		Result: &jsonAdaptor,
		URL: rest.URL{
			Path: kubeAPIIngressRootPath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
	return ingressList, err
}

// GetIngressList -- returns ingress with specified domain.
func (client *Client) GetIngress(namespace, domain string) (model.Ingress, error) {
	var ingress model.Ingress
	err := client.RestAPI.Get(rest.Rq{
		Result: &ingress,
		URL: rest.URL{
			Path: kubeAPIIngressPath,
			Params: rest.P{
				"namespace": namespace,
				"domain":    domain,
			},
		},
	})
	return ingress, err
}

// UpdateIngress -- updates ingress on provided domain with new one
func (client *Client) UpdateIngress(namespace, domain string, ingress model.Ingress) error {
	return client.RestAPI.Put(rest.Rq{
		Body: ingress,
		URL: rest.URL{
			Path: resourceIngressPath,
			Params: rest.P{
				"namespace": namespace,
				"domain":    domain,
			},
		},
	})
}

// DeleteIngress -- deletes ingress on provided domain
func (client *Client) DeleteIngress(namespace, domain string) error {
	return client.RestAPI.Put(rest.Rq{
		URL: rest.URL{
			Path: resourceIngressPath,
			Params: rest.P{
				"namespace": namespace,
				"domain":    domain,
			},
		},
	})
}
