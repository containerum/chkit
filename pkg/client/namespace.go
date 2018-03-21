package chClient

import (
	"fmt"

	"git.containerum.net/ch/kube-client/pkg/cherry"
	"git.containerum.net/ch/kube-client/pkg/cherry/auth"
	"git.containerum.net/ch/kube-client/pkg/cherry/kube-api"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model/namespace"
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
	err := retry(4, func() error {
		kubeNamespace, err := client.kubeAPIClient.GetNamespace(label)
		switch {
		case err == nil:
			ns = namespace.NamespaceFromKube(kubeNamespace)
			return nil
		case cherry.In(err,
			kubeErrors.ErrResourceNotExist(),
			kubeErrors.ErrAccessError(),
			kubeErrors.ErrUnableGetResource()):
			return ErrNamespaceNotExists
		case cherry.In(err, autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			return client.Auth()
		default:
			return err
		}
	})
	return ns, err
}

func (client *Client) GetNamespaceList() (namespace.NamespaceList, error) {
	var list namespace.NamespaceList
	err := retry(4, func() error {
		kubeList, err := client.kubeAPIClient.GetNamespaceList(nil)
		switch {
		case err == nil:
			list = namespace.NamespaceListFromKube(kubeList)
			return err
		case cherry.Equals(err, autherr.ErrInvalidToken()) ||
			cherry.Equals(err, autherr.ErrTokenNotFound()):
			fmt.Printf("reauth: %v\n", err)
			return client.Auth()
		case cherry.Equals(err, kubeErrors.ErrAccessError()):
			return ErrYouDoNotHaveAccessToNamespace
		default:
			return err
		}
	})
	return list, err
}
