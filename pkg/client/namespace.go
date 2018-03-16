package chClient

import (
	"git.containerum.net/ch/kube-client/pkg/cherry"
	"git.containerum.net/ch/kube-client/pkg/cherry/auth"
	"git.containerum.net/ch/kube-client/pkg/cherry/resource-service"
	kubeClientModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"
)

const (
	// ErrUnableToGetNamespace -- unable to get namespace
	ErrUnableToGetNamespace chkitErrors.Err = "unable to get namespace"
	ErrNamespaceNotExists   chkitErrors.Err = "namespace not exists"
)

// GetNamespace -- returns info of namespace with given label
func (client *Client) GetNamespace(label string) (model.Namespace, error) {
	var err error
	var namespace kubeClientModels.Namespace
	for i := 0; i == 0 || (i < 4 && err != nil); i++ {
		namespace, err = client.kubeAPIClient.GetNamespace(label)
		switch {
		case err == nil:
			break
		case cherry.Equals(err, rserrors.ErrResourceNotExists()):
			return model.Namespace{}, ErrNamespaceNotExists
		case cherry.Equals(err, autherr.ErrInvalidToken()) ||
			cherry.Equals(err, autherr.ErrTokenNotFound()):
			switch client.Auth() {
			case ErrWrongPasswordLoginCombination, ErrUserNotExist:
				return model.Namespace{}, err
			}
		}
	}
	return model.NamespaceFromKube(namespace), err
}
