package client

import (
	permModel "git.containerum.net/ch/permissions/pkg/model"
	"github.com/containerum/kube-client/pkg/rest"
)

const (
	//TODO: Change path to Permissions "/namespaces/{namespace}/accesses"
	accessesPath = "/namespace/{namespace}/access"
)

func (client *Client) GetNamespaceAccesses(namespace string) (permModel.NamespaceWithPermissionsJSON, error) {
	var access permModel.NamespaceWithPermissionsJSON
	err := client.RestAPI.Get(rest.Rq{
		Result: &access,
		URL: rest.URL{
			Path: accessesPath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
	return access, err
}
