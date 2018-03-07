package client

import (
	"strconv"

	"git.containerum.net/ch/kube-client/pkg/model"
	"git.containerum.net/ch/kube-client/pkg/rest"
)

const (
	resourceIngressRootPath = "/namespace/{namespace}/ingress"
	resourceIngressPath     = resourceIngressRootPath + "/{domain}"
)

// AddIngress -- adds ingress to provided namespace
func (client *Client) AddIngress(namespace string, ingress model.Ingress) error {
	return client.RestAPI.Post(rest.Rq{
		Body: ingress,
		URL: rest.URL{
			Path: client.APIurl + resourceIngressRootPath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
}

// GetIngressList -- returns list of ingresses.
// Consumes namespace and optional pagination parameters
// If role=admin && !user-id -> return all
// If role=admin && user-id -> return user's
// If role=user -> return user's
func (client *Client) GetIngressList(namespace string, page, perPage *uint64) ([]model.Ingress, error) {
	var ingressList []model.Ingress
	err := client.RestAPI.Get(rest.Rq{
		Result: &ingressList,
		Query: rest.Q{
			"page":     strconv.FormatUint(*page, 10),
			"per_page": strconv.FormatUint(*perPage, 10),
		},
		URL: rest.URL{
			Path: client.APIurl + resourceIngressRootPath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
	return ingressList, err
}

// UpdateIngress -- updates ingress on provided domain with new one
func (client *Client) UpdateIngress(namespace, domain string, ingress model.Ingress) error {
	return client.RestAPI.Put(rest.Rq{
		Body: ingress,
		URL: rest.URL{
			Path: client.APIurl + resourceIngressPath,
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
			Path: client.APIurl + resourceIngressPath,
			Params: rest.P{
				"namespace": namespace,
				"domain":    domain,
			},
		},
	})
}
