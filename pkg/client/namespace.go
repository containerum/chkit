package chClient

import (
	"fmt"
	"net/http"

	"git.containerum.net/ch/kube-client/pkg/cherry"

	kubeClientModels "git.containerum.net/ch/kube-client/pkg/model"
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
	var err error
	var namespace kubeClientModels.Namespace
	for i := 0; i == 0 || (i < 4 && err != nil); i++ {
		namespace, err = client.kubeAPIClient.GetNamespace(label)
		switch er := err.(type) {
		case nil:
			break
		case *rest.UnexpectedHTTPstatusError:
			if er.Status == http.StatusNotFound {
				return model.Namespace{}, ErrUnableToGetNamespace.
					Wrap(fmt.Errorf("namespace %q doesn't exist", label))
			}
		case *cherry.Err:
			switch er.ID.Kind {
			case 2, 3:
				// rotten access token
				if err = client.Auth(); err != nil {
					continue
				}
			}
		}
	}
	return model.NamespaceFromKube(namespace), err
}
