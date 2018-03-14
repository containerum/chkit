package chClient

import (
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"
)

const (
	// ErrUnableToGetNamespace -- unable to get namespace
	ErrUnableToGetNamespace chkitErrors.Err = "unable to get namespace"
)

// GetNamespace -- returns info of namespace with given label
func (client *Client) GetNamespace(label string) (model.Namespace, error) {
	namespace, err := client.kubeAPIClient.GetNamespace(label)
	if err != nil {
		return model.Namespace{}, ErrUnableToGetNamespace.Wrap(err)
	}
	return model.NamespaceFromKube(namespace), nil
}
