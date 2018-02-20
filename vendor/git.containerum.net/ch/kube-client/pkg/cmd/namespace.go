package cmd

import (
	"fmt"
	"net/http"
	"strconv"

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
		Get(client.APIurl + getNamespaceList)
	if err := catchErr(err, resp, http.StatusOK); err != nil {
		return []model.Namespace{}, err
	}
	return *resp.Result().(*[]model.Namespace), nil
}

//GetNamespace return namespace by Name
func (client *Client) GetNamespace(ns string) (model.Namespace, error) {
	resp, err := client.Request.SetResult(model.Namespace{}).
		SetPathParams(map[string]string{
			"namespace": ns,
		}).
		Get(client.APIurl + getNamespace)
	if err := catchErr(err, resp, http.StatusOK); err != nil {
		return model.Namespace{}, err
	}
	return *resp.Result().(*model.Namespace), nil
}

// ResourceGetNamespace -- consumes a namespace and an optional user ID
// returns a namespace data OR an error
func (client *Client) ResourceGetNamespace(namespace string, userID *string) (model.ResourceNamespace, error) {
	req := client.Request.
		SetPathParams(map[string]string{
			"namespace": namespace,
		}).SetResult(model.ResourceNamespace{}).
		SetError(model.ResourceError{})
	if userID != nil {
		req.SetQueryParam("user-id", *userID)
	}
	resp, err := req.Get(client.ResourceAddr + resourceNamespacePath)
	if err := catchErr(err, resp, http.StatusOK); err != nil {
		return model.ResourceNamespace{}, err
	}
	return *resp.Result().(*model.ResourceNamespace), nil
}

// ResourceGetNamespaceList -- consumes a page number parameter,
// amount of namespaces per page and optional userID,
// returns a slice of Namespaces OR a nil slice AND an error
func (client *Client) ResourceGetNamespaceList(page, perPage uint64, userID string) ([]model.ResourceNamespace, error) {
	req := client.Request.
		SetQueryParams(map[string]string{
			"page":     strconv.FormatUint(page, 10),
			"per_page": strconv.FormatUint(perPage, 10),
		}).SetResult([]model.ResourceNamespace{}).
		SetError(model.ResourceError{})
	if userID != "" {
		req.SetQueryParam("user-id", userID)
	}
	resp, err := req.Get(client.ResourceAddr + resourceNamespacesPath)
	if err := catchErr(err, resp, http.StatusOK); err != nil {
		return nil, err
	}
	return *resp.Result().(*[]model.ResourceNamespace), nil
}

// RenameNamespace -- renames user namespace
// Consumes namespace name and new name.
func (client *Client) RenameNamespace(namespace, newName string) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace": resourceNamespacePath,
		}).SetBody(model.ResourceUpdateName{
		Label: newName,
	}).Put(client.ResourceAddr + resourceNamespacePath)
	if err != nil {
		return err
	}
	switch resp.StatusCode() {
	case http.StatusOK, http.StatusAccepted:
		return nil
	default:
		if resp.Error() != nil {
			return fmt.Errorf("%v", resp.Error())
		}
		return fmt.Errorf("%v", resp.Status())
	}
}

// SetNamespaceAccess -- sets/changes access to namespace for provided user
func (client *Client) SetNamespaceAccess(namespace, username, access string) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace": namespace,
		}).SetBody(model.ResourceUpdateUserAccess{
		Username: username,
		Access:   access,
	}).Post(client.ResourceAddr + resourceNamespaceNamePath)
	if err != nil {
		return err
	}
	switch resp.StatusCode() {
	case http.StatusOK, http.StatusAccepted:
		return nil
	default:
		if resp.Error() != nil {
			return fmt.Errorf("%v", resp.Error())
		}
		return fmt.Errorf("%v", resp.Status())
	}
}

// DeleteNamespaceAccess -- deletes user access to namespace
func (client *Client) DeleteNamespaceAccess(namespace, username string) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace": namespace,
		}).SetBody(model.ResourceUpdateUserAccess{
		Username: username,
	}).Delete(client.ResourceAddr + resourceNamespaceNamePath)
	if err != nil {
		return err
	}
	switch resp.StatusCode() {
	case http.StatusOK, http.StatusAccepted:
		return nil
	default:
		if resp.Error() != nil {
			return fmt.Errorf("%v", resp.Error())
		}
		return fmt.Errorf("%v", resp.Status())
	}
}

// DeleteNamespace -- deletes provided namespace
func (client *Client) DeleteNamespace(namespace string) error {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"namespace": namespace,
		}).SetError(model.ResourceError{}).
		Delete(client.ResourceAddr + getNamespace)
	return catchErr(err, resp, http.StatusOK, http.StatusAccepted)
}
