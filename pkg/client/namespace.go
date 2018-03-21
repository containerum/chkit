package chClient

import (
	"fmt"

	"git.containerum.net/ch/kube-client/pkg/cherry"
	"git.containerum.net/ch/kube-client/pkg/cherry/auth"
	"git.containerum.net/ch/kube-client/pkg/cherry/kube-api"
	kubeClientModels "git.containerum.net/ch/kube-client/pkg/model"
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
	var err error
	var kubeNamespace kubeClientModels.Namespace
retry:
	for i := uint(0); i == 0 || (i < 4 && err != nil); i++ {
		kubeNamespace, err = client.kubeAPIClient.GetNamespace(label)
		switch {
		case err == nil:
			break retry
		case cherry.In(err,
			kubeErrors.ErrResourceNotExist(),
			kubeErrors.ErrAccessError(),
			kubeErrors.ErrUnableGetResource()):
			return namespace.Namespace{}, ErrNamespaceNotExists
		case cherry.In(err, autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			if er := client.Auth(); er != nil {
				return namespace.Namespace{}, er
			}
		}
		waitNextAttempt(i)
	}
	return namespace.NamespaceFromKube(kubeNamespace), err
}

func (client *Client) GetNamespaceList() (namespace.NamespaceList, error) {
	var err error
	var list []kubeClientModels.Namespace
	for i := uint(0); i == 0 || (i < 4 && err != nil); i++ {
		list, err = client.kubeAPIClient.GetNamespaceList(nil)
		switch {
		case err == nil:
			break
		case cherry.Equals(err, autherr.ErrInvalidToken()) ||
			cherry.Equals(err, autherr.ErrTokenNotFound()):
			fmt.Printf("reauth: %v\n", err)
			err = client.Auth()
			switch err {
			case ErrWrongPasswordLoginCombination, ErrUserNotExist:
				return []namespace.Namespace{}, err
			default:
				fmt.Println(err)
			}
		case cherry.Equals(err, kubeErrors.ErrAccessError()):
			return namespace.NamespaceList{}, ErrYouDoNotHaveAccessToNamespace
		}
		waitNextAttempt(i)
	}
	return namespace.NamespaceListFromKube(list), err
}
