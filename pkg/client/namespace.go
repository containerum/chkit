package chClient

import (
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"
)

const (
	ErrUnableToGetNamespace chkitErrors.Err = "unable to get namespace"
)

func (client *Client) GetNamespace(label string) (model.Namespace, error) {
	namespace, err := client.kubeApiClient.GetNamespace(label)
	if err != nil {
		return model.Namespace{}, ErrUnableToGetNamespace.Wrap(err)
	}
	return model.NamespaceFromKube(namespace), nil
}
