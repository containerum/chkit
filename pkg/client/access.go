package chClient

import (
	"git.containerum.net/ch/auth/pkg/errors"
	"git.containerum.net/ch/kube-api/pkg/kubeErrors"
	"github.com/containerum/cherry"
	"github.com/containerum/chkit/pkg/model/access"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/sirupsen/logrus"
)

func (client *Client) GetAccess(nsName string) (access.AccessList, error) {
	nswp, err := client.kubeAPIClient.GetNamespaceAccesses(nsName)
	if err != nil {
		return nil, err
	}
	return access.AccessListFromKube(nswp), err
}

/*func (client *Client) GetAccessList() (access.AccessList, error) {
	list, err := client.GetNamespaceList()
	return access.AccessListFromNamespaces(list), err
}*/

func (client *Client) SetAccess(ns, username string, acc model.AccessLevel) error {
	err := retry(4, func() (bool, error) {
		err := client.kubeAPIClient.SetNamespaceAccess(ns, username, acc)
		switch {
		case err == nil:
			return false, nil
		case cherry.In(err,
			kubeErrors.ErrResourceNotExist(),
			kubeErrors.ErrAccessError(),
			kubeErrors.ErrUnableGetResource()):
			return false, err
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound(),
			autherr.ErrTokenNotOwnedBySender()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	if err != nil {
		logrus.WithError(err).WithField("namespace", ns).
			Errorf("unable to set access to namespace")
	}
	return err
}

func (client *Client) DeleteAccess(ns, username string) error {
	err := retry(4, func() (bool, error) {
		err := client.kubeAPIClient.DeleteNamespaceAccess(ns, username)
		switch {
		case err == nil:
			return false, nil
		case cherry.In(err,
			kubeErrors.ErrResourceNotExist(),
			kubeErrors.ErrAccessError(),
			kubeErrors.ErrUnableGetResource()):
			return false, err
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound(),
			autherr.ErrTokenNotOwnedBySender()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	if err != nil {
		logrus.WithError(err).WithField("namespace", ns).
			Errorf("unable to delete access to namespace")
	}
	return err
}
