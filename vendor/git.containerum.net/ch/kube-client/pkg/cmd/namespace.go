package cmd

import (
	"git.containerum.net/ch/kube-client/pkg/model"
)

//ListOptions -
type ListOptions struct {
	Owner string
}

const (
	getNamespace     = "/namespaces/{namespace}"
	getNamespaceList = "/namespaces"
)

//GetNamespaceList return namespace list. Can use query filters: owner
func (c *Client) GetNamespaceList(queries map[string]string) ([]model.Namespace, error) {
	resp, err := c.Request.
		SetQueryParams(queries).
		SetResult([]model.Namespace{}).
		Get(c.serverURL + getNamespaceList)
	if err != nil {
		return []model.Namespace{}, err
	}
	return *resp.Result().(*[]model.Namespace), nil
}

//GetNamespace return namespace by Name
func (c *Client) GetNamespace(ns string) (model.Namespace, error) {
	resp, err := c.Request.SetResult(model.Namespace{}).
		SetPathParams(map[string]string{
			"namespace": ns,
		}).
		Get(c.serverURL + getNamespace)
	if err != nil {
		return model.Namespace{}, err
	}
	return *resp.Result().(*model.Namespace), nil
}
