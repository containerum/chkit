package chClient

import (
	"git.containerum.net/ch/kube-client/pkg/cherry"
	"git.containerum.net/ch/kube-client/pkg/cherry/auth"
	"git.containerum.net/ch/kube-client/pkg/cherry/kube-api"
	"git.containerum.net/ch/kube-client/pkg/cherry/resource-service"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model/namespace"
	"github.com/sirupsen/logrus"
)

const (
	// ErrUnableToGetNamespace -- unable to get namespace
	ErrUnableToGetNamespace chkitErrors.Err = "unable to get namespace"
	// ErrYouDoNotHaveAccessToNamespace -- you don't have access to namespace
	ErrYouDoNotHaveAccessToNamespace chkitErrors.Err = "you don't have access to namespace"
	// ErrNamespaceNotExists -- namespace not exists
	ErrNamespaceNotExists chkitErrors.Err = "namespace not exists"
)

// GetNamespace -- returns info of namespace with given label.
// Returns:
// 	- ErrNamespaceNotExists
//  - ErrWrongPasswordLoginCombination
//  - ErrUserNotExist
func (client *Client) GetNamespace(label string) (namespace.Namespace, error) {
	var ns namespace.Namespace
	err := retry(4, func() (bool, error) {
		kubeNamespace, err := client.kubeAPIClient.GetNamespace(label)
		switch {
		case err == nil:
			ns = namespace.NamespaceFromKube(kubeNamespace)
			return false, nil
		case cherry.In(err,
			kubeErrors.ErrResourceNotExist(),
			kubeErrors.ErrAccessError(),
			kubeErrors.ErrUnableGetResource()):
			return false, ErrNamespaceNotExists
		case cherry.In(err, autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	return ns, err
}

func (client *Client) GetNamespaceList() (namespace.NamespaceList, error) {
	var list namespace.NamespaceList
	err := retry(4, func() (bool, error) {
		kubeList, err := client.kubeAPIClient.GetNamespaceList(nil)
		switch {
		case err == nil:
			list = namespace.NamespaceListFromKube(kubeList)
			return false, err
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			er := client.Auth()
			return true, er
		case cherry.In(err, kubeErrors.ErrAccessError()):
			return false, ErrYouDoNotHaveAccessToResource.Wrap(err)
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	return list, err
}

func (client *Client) DeleteNamespace(namespace string) error {
	return retry(4, func() (bool, error) {
		err := client.kubeAPIClient.DeleteNamespace(namespace)
		switch {
		case err == nil:
			return false, nil
		case cherry.In(err,
			rserrors.ErrResourceNotExists()):
			logrus.WithError(ErrResourceNotExists.Wrap(err)).
				Debugf("error while deleting namespace %q", namespace)
			return false, ErrResourceNotExists
		case cherry.In(err,
			rserrors.ErrResourceNotOwned(),
			rserrors.ErrAccessRecordNotExists()):
			logrus.WithError(ErrYouDoNotHaveAccessToResource.Wrap(err)).
				Debugf("error while deleting namespace %q", namespace)
			return false, ErrYouDoNotHaveAccessToResource
		case cherry.In(err, autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			err = client.Auth()
			if err != nil {
				logrus.WithError(err).
					Debugf("error while deleting namespace %q", namespace)
			}
			return true, err
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
}
