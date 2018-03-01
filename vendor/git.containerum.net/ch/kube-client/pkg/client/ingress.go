package client

import (
	"net/http"
	"strconv"

	"git.containerum.net/ch/kube-client/pkg/cherry"
	"git.containerum.net/ch/kube-client/pkg/model"
)

const (
	resourceIngressRootPath = "/namespace/{namespace}/ingress"
	resourceIngressPath     = resourceIngressRootPath + "/{domain}"
)

// AddIngress -- adds ingress to provided namespace
func (client *Client) AddIngress(namespace string, ingress model.Ingress) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace": namespace,
		}).SetBody(ingress).
		SetError(cherry.Err{}).
		Post(client.ResourceAddr + resourceIngressRootPath)
	return MapErrors(resp, err,
		http.StatusOK,
		http.StatusAccepted)
}

// GetIngressList -- returns list of ingresses.
// Consumes namespace and optional pagination parameters
// If role=admin && !user-id -> return all
// If role=admin && user-id -> return user's
// If role=user -> return user's
func (client *Client) GetIngressList(namespace string, page, perPage *uint64) ([]model.Ingress, error) {
	req := client.Request.
		SetPathParams(map[string]string{
			"namespace": namespace,
		}).
		SetResult([]model.Ingress{}).
		SetError(cherry.Err{})
	if page != nil {
		req.SetQueryParam("page", strconv.FormatUint(*page, 10))
	}
	if perPage != nil {
		req.SetQueryParam("per_page", strconv.FormatUint(*perPage, 10))
	}
	resp, err := req.Get(client.ResourceAddr + resourceIngressRootPath)
	if err = MapErrors(resp, err, http.StatusOK); err != nil {
		return nil, err
	}
	return *resp.Result().(*[]model.Ingress), nil

}

// UpdateIngress -- updates ingress on provided domain with new one
func (client *Client) UpdateIngress(namespace, domain string, ingress model.Ingress) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace": namespace,
			"domain":    domain,
		}).SetBody(ingress).
		SetError(cherry.Err{}).
		Put(client.ResourceAddr + resourceIngressPath)
	return MapErrors(resp, err,
		http.StatusOK,
		http.StatusAccepted)
}

// DeleteIngress -- deletes ingress on provided domain
func (client *Client) DeleteIngress(namespace, domain string) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace": namespace,
			"domain":    domain,
		}).
		SetError(cherry.Err{}).
		Delete(client.ResourceAddr + resourceIngressPath)
	return MapErrors(resp, err,
		http.StatusOK,
		http.StatusAccepted,
		http.StatusNoContent)
}
