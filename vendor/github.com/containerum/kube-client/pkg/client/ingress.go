package client

import (
	"github.com/containerum/kube-client/pkg/model"
	"github.com/containerum/kube-client/pkg/rest"
)

const (
	ingressesPath = "/namespaces/{namespace}/ingresses"
	ingressPath   = "/namespaces/{namespace}/ingresses/{domain}"
)

// AddIngress -- adds ingress to provided namespace
func (client *Client) AddIngress(namespace string, ingress model.Ingress) error {
	return client.RestAPI.Post(rest.Rq{
		Body: ingress,
		URL: rest.URL{
			Path: ingressesPath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
}

// GetIngressList -- returns list of ingresses.
func (client *Client) GetIngressList(namespace string) (model.IngressesList, error) {
	var ingressList model.IngressesList
	err := client.RestAPI.Get(rest.Rq{
		Result: &ingressList,
		URL: rest.URL{
			Path: ingressesPath,
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
			Path: ingressPath,
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
			Path: ingressPath,
			Params: rest.P{
				"namespace": namespace,
				"domain":    domain,
			},
		},
	})
}

// DeleteIngress -- deletes ingress on provided domain
func (client *Client) DeleteIngress(namespace, domain string) error {
	return client.RestAPI.Delete(rest.Rq{
		URL: rest.URL{
			Path: ingressPath,
			Params: rest.P{
				"namespace": namespace,
				"domain":    domain,
			},
		},
	})
}
