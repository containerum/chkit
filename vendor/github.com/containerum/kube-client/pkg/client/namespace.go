package client

import (
	"strconv"

	"github.com/containerum/kube-client/pkg/model"
	"github.com/containerum/kube-client/pkg/rest"
)

const (
	namespacesPath      = "/namespaces"
	namespacePath       = "/namespaces/{namespace}"
	namespaceNamePath   = "/namespaces/{namespace}/rename"
	namespaceAccessPath = "/namespaces/{namespace}/accesses"
)

//GetNamespaceList return namespace list. Can use query filters: owner
func (client *Client) GetNamespaceList() (model.NamespacesList, error) {
	var namespaceList model.NamespacesList
	err := client.RestAPI.Get(rest.Rq{
		Result: &namespaceList,
		URL: rest.URL{
			Path:   namespacesPath,
			Params: rest.P{},
		},
	})
	return namespaceList, err
}

//GetNamespace return namespace by ID
func (client *Client) GetNamespace(ns string) (model.Namespace, error) {
	var namespace model.Namespace
	err := client.RestAPI.Get(rest.Rq{
		Result: &namespace,
		URL: rest.URL{
			Path: namespacePath,
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
			Path: namespacePath,
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
func (client *Client) ResourceGetNamespaceList(page, perPage uint64) (model.NamespacesList, error) {
	var namespaceList model.NamespacesList
	err := client.RestAPI.Get(rest.Rq{
		Result: &namespaceList,
		Query: rest.Q{
			"page":     strconv.FormatUint(page, 10),
			"per_page": strconv.FormatUint(perPage, 10),
		},
		URL: rest.URL{
			Path:   namespacesPath,
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
			Path: namespaceNamePath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
}

// SetNamespaceAccess -- sets/changes access to namespace for provided user
func (client *Client) SetNamespaceAccess(namespace, username string, access model.AccessLevel) error {
	return client.RestAPI.Put(rest.Rq{
		Body: model.ResourceUpdateUserAccess{
			Username: username,
			Access:   access,
		},
		URL: rest.URL{
			Path: namespaceAccessPath,
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
			Path: namespaceAccessPath,
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
			Path: namespacePath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
}
