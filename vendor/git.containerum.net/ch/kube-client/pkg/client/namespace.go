package client

import (
	"strconv"

	"git.containerum.net/ch/kube-client/pkg/rest"

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
	var namespaceList []model.Namespace
	err := client.RestAPI.Get(rest.Rq{
		Result: &namespaceList,
		Query:  queries,
		URL: rest.URL{
			Path:   client.APIurl + getNamespaceList,
			Params: rest.P{},
		},
	})
	return namespaceList, err
}

//GetNamespace return namespace by Name
func (client *Client) GetNamespace(ns string) (model.Namespace, error) {
	var namespace model.Namespace
	err := client.RestAPI.Get(rest.Rq{
		Result: &namespace,
		URL: rest.URL{
			Path: client.APIurl + getNamespace,
			Params: rest.P{
				"namespace": ns,
			},
		},
	})
	return namespace, err
}

// ResourceGetNamespace -- consumes a namespace
// returns a namespace data OR an error
func (client *Client) ResourceGetNamespace(namespace string) (model.Namespace, error) {
	var ns model.Namespace
	err := client.RestAPI.Get(rest.Rq{
		Result: &ns,
		URL: rest.URL{
			Path: client.APIurl + resourceNamespacePath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
	return ns, err
}

// ResourceGetNamespaceList -- consumes a page number parameter,
// amount of namespaces per page and optional userID,
// returns a slice of Namespaces OR a nil slice AND an error
func (client *Client) ResourceGetNamespaceList(page, perPage uint64) ([]model.Namespace, error) {
	var namespaceList []model.Namespace
	err := client.RestAPI.Get(rest.Rq{
		Result: &namespaceList,
		Query: rest.Q{
			"page":     strconv.FormatUint(page, 10),
			"per_page": strconv.FormatUint(perPage, 10),
		},
		URL: rest.URL{
			Path:   client.APIurl + resourceNamespacesPath,
			Params: rest.P{},
		},
	})
	return namespaceList, err
}

// RenameNamespace -- renames user namespace
// Consumes namespace name and new name.
func (client *Client) RenameNamespace(namespace, newName string) error {
	return client.RestAPI.Put(rest.Rq{
		Body: model.ResourceUpdateName{
			Label: newName,
		},
		URL: rest.URL{
			Path: client.APIurl + resourceNamespacePath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
}

// SetNamespaceAccess -- sets/changes access to namespace for provided user
func (client *Client) SetNamespaceAccess(namespace, username, access string) error {
	return client.RestAPI.Post(rest.Rq{
		Body: model.ResourceUpdateUserAccess{
			Username: username,
			Access:   access,
		},
		URL: rest.URL{
			Path: client.APIurl + resourceNamespaceNamePath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
}

// DeleteNamespaceAccess -- deletes user access to namespace
func (client *Client) DeleteNamespaceAccess(namespace, username string) error {
	return client.RestAPI.Delete(rest.Rq{
		Body: model.ResourceUpdateUserAccess{
			Username: username,
		},
		URL: rest.URL{
			Path: client.APIurl + resourceNamespaceNamePath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
}

// DeleteNamespace -- deletes provided namespace
func (client *Client) DeleteNamespace(namespace string) error {
	return client.RestAPI.Delete(rest.Rq{
		URL: rest.URL{
			Path: client.APIurl + getNamespace,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
}
