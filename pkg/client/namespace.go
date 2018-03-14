package chClient

import (
	"fmt"

	"git.containerum.net/ch/kube-client/pkg/rest"
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
		switch err := err.(type) {
		case *rest.UnexpectedHTTPstatusError:
			if err.Status == 404 {
				return model.Namespace{}, ErrUnableToGetNamespace.
					Wrap(fmt.Errorf("namespace %q doesn't exist", label))
			}
		}
		return model.Namespace{}, ErrUnableToGetNamespace.Wrap(err)
	}
	return model.NamespaceFromKube(namespace), nil
}
