package chClient

import (
	"git.containerum.net/ch/auth/pkg/errors"
	"git.containerum.net/ch/kube-api/pkg/kubeErrors"
	permErrors "git.containerum.net/ch/permissions/pkg/errors"
	"git.containerum.net/ch/resource-service/pkg/rsErrors"
	"github.com/containerum/cherry"
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
func (client *Client) GetNamespace(ID string) (namespace.Namespace, error) {
	var ns namespace.Namespace
	err := retry(4, func() (bool, error) {
		kubeNamespace, err := client.kubeAPIClient.GetNamespace(ID)
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
		kubeList, err := client.kubeAPIClient.GetNamespaceList()
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

func (client *Client) DeleteNamespace(ID string) error {
	return retry(4, func() (bool, error) {
		err := client.kubeAPIClient.DeleteNamespace(ID)
		switch {
		case err == nil:
			return false, nil
		case cherry.In(err,
			rserrors.ErrResourceNotExists()):
			logrus.WithError(ErrResourceNotExists.Wrap(err)).
				Debugf("error while deleting ID %q", ID)
			return false, ErrResourceNotExists
		case cherry.In(err,
			permErrors.ErrResourceNotOwned(),
			rserrors.ErrPermissionDenied()):
			logrus.WithError(ErrYouDoNotHaveAccessToResource.Wrap(err)).
				Debugf("error while deleting ID %q", ID)
			return false, ErrYouDoNotHaveAccessToResource
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			err = client.Auth()
			if err != nil {
				logrus.WithError(err).
					Debugf("error while deleting ID %q", ID)
			}
			return true, err
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
}

func (client *Client) RenameNamespace(ID, newName string) error {
	err := retry(4, func() (bool, error) {
		err := client.kubeAPIClient.RenameNamespace(ID, newName)
		switch {
		case err == nil:
			return false, nil
		case cherry.In(err,
			rserrors.ErrResourceNotExists()):
			return false, ErrResourceNotExists.Wrap(err)
		case cherry.In(err,
			permErrors.ErrResourceNotOwned(),
			rserrors.ErrPermissionDenied()):
			return false, ErrYouDoNotHaveAccessToResource.Wrap(err)
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	if err != nil {
		logrus.WithError(err).Errorf("unable to rename namespace %q", ID)
	}
	return err
}
