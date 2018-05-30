package client

import (
	"github.com/containerum/kube-client/pkg/model"
	"github.com/containerum/kube-client/pkg/rest"
)

const (
	accessesPath = "/namespaces/{namespace}/access"
)

func (client *Client) GetNamespaceAccesses(namespace string) ([]model.UserAccess, error) {
	var access model.Namespace
	err := client.RestAPI.Get(rest.Rq{
		Result: &access,
		URL: rest.URL{
			Path: accessesPath,
			Params: rest.P{
				"namespace": namespace,
			},
		},
	})
	return access.Users, err
}
