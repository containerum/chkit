package client

import (
	"net/http"
	"strconv"

	"git.containerum.net/ch/kube-client/pkg/cherry"

	"git.containerum.net/ch/kube-client/pkg/model"
)

//ListOptions -
type ListOptions struct {
	Owner string
}

const (
	getNamespace                = "/namespaces/{namespace}"
	getNamespaceList            = "/namespaces"
	resourceNamespacePath       = "/namespace/{namespace}"
	resourceNamespacesPath      = "/namespace"
	resourceNamespaceNamePath   = resourceNamespacePath + "/name"
	resourceNamespaceAccessPath = resourceNamespacePath + "/access"
)

//GetNamespaceList return namespace list. Can use query filters: owner
func (client *Client) GetNamespaceList(queries map[string]string) ([]model.Namespace, error) {
	resp, err := client.Request.
		SetQueryParams(queries).
		SetResult([]model.Namespace{}).
		SetError(cherry.Err{}).
		Get(client.APIurl + getNamespaceList)
	if err := MapErrors(resp, err, http.StatusOK); err != nil {
		return nil, err
	}
	return *resp.Result().(*[]model.Namespace), nil
}

//GetNamespace return namespace by Name
func (client *Client) GetNamespace(ns string) (model.Namespace, error) {
	resp, err := client.Request.
		SetResult(model.Namespace{}).
		SetPathParams(map[string]string{
			"namespace": ns,
		}).
		SetError(cherry.Err{}).
		Get(client.APIurl + getNamespace)
	if err := MapErrors(resp, err, http.StatusOK); err != nil {
		return model.Namespace{}, err
	}
	return *resp.Result().(*model.Namespace), nil
}

// ResourceGetNamespace -- consumes a namespace
// returns a namespace data OR an error
func (client *Client) ResourceGetNamespace(namespace string) (model.Namespace, error) {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace": namespace,
		}).SetResult(model.Namespace{}).
		SetError(cherry.Err{}).
		Get(client.ResourceAddr + resourceNamespacePath)
	if err := MapErrors(resp, err, http.StatusOK); err != nil {
		return model.Namespace{}, err
	}
	return *resp.Result().(*model.Namespace), nil
}

// ResourceGetNamespaceList -- consumes a page number parameter,
// amount of namespaces per page and optional userID,
// returns a slice of Namespaces OR a nil slice AND an error
func (client *Client) ResourceGetNamespaceList(page, perPage uint64) ([]model.Namespace, error) {
	req := client.Request.
		SetQueryParams(map[string]string{
			"page":     strconv.FormatUint(page, 10),
			"per_page": strconv.FormatUint(perPage, 10),
		}).SetResult([]model.Namespace{}).
		SetError(cherry.Err{})
	resp, err := req.Get(client.ResourceAddr + resourceNamespacesPath)
	if err := MapErrors(resp, err, http.StatusOK); err != nil {
		return nil, err
	}
	return *resp.Result().(*[]model.Namespace), nil
}

// RenameNamespace -- renames user namespace
// Consumes namespace name and new name.
func (client *Client) RenameNamespace(namespace, newName string) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace": namespace,
		}).SetBody(model.ResourceUpdateName{
		Label: newName,
	}).SetError(cherry.Err{}).
		Put(client.ResourceAddr + resourceNamespacePath)
	return MapErrors(resp, err,
		http.StatusOK,
		http.StatusAccepted)
}

// SetNamespaceAccess -- sets/changes access to namespace for provided user
func (client *Client) SetNamespaceAccess(namespace, username, access string) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace": namespace,
		}).SetBody(model.ResourceUpdateUserAccess{
		Username: username,
		Access:   access,
	}).SetError(cherry.Err{}).
		Post(client.ResourceAddr + resourceNamespaceNamePath)
	return MapErrors(resp, err,
		http.StatusOK,
		http.StatusAccepted)
}

// DeleteNamespaceAccess -- deletes user access to namespace
func (client *Client) DeleteNamespaceAccess(namespace, username string) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace": namespace,
		}).SetBody(model.ResourceUpdateUserAccess{
		Username: username,
	}).SetError(cherry.Err{}).
		Delete(client.ResourceAddr + resourceNamespaceNamePath)
	return MapErrors(resp, err,
		http.StatusOK,
		http.StatusAccepted)
}

// DeleteNamespace -- deletes provided namespace
func (client *Client) DeleteNamespace(namespace string) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace": namespace,
		}).SetError(cherry.Err{}).
		Delete(client.ResourceAddr + getNamespace)
	return MapErrors(resp, err,
		http.StatusOK,
		http.StatusAccepted)
}
