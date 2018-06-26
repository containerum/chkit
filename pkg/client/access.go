package chClient

import (
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
		return HandleErrorRetry(client, err)
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
		return HandleErrorRetry(client, err)
	})
	if err != nil {
		logrus.WithError(err).WithField("namespace", ns).
			Errorf("unable to delete access to namespace")
	}
	return err
}
