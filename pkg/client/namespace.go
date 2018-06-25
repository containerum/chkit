package chClient

import (
	"github.com/containerum/chkit/pkg/model/namespace"
	"github.com/sirupsen/logrus"
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
		if err == nil {
			ns = namespace.NamespaceFromKube(kubeNamespace)
		}
		return HandleErrorRetry(client, err)
	})
	return ns, err
}

func (client *Client) GetNamespaceList() (namespace.NamespaceList, error) {
	var list namespace.NamespaceList
	err := retry(4, func() (bool, error) {
		kubeList, err := client.kubeAPIClient.GetNamespaceList()
		if err == nil {
			list = namespace.NamespaceListFromKube(kubeList)
		}
		return HandleErrorRetry(client, err)
	})
	return list, err
}

func (client *Client) DeleteNamespace(ID string) error {
	return retry(4, func() (bool, error) {
		err := client.kubeAPIClient.DeleteNamespace(ID)
		return HandleErrorRetry(client, err)
	})
}

func (client *Client) RenameNamespace(ID, newName string) error {
	err := retry(4, func() (bool, error) {
		err := client.kubeAPIClient.RenameNamespace(ID, newName)
		return HandleErrorRetry(client, err)
	})
	if err != nil {
		logrus.WithError(err).Errorf("unable to rename namespace %q", ID)
	}
	return err
}
