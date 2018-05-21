package chClient

import (
	"git.containerum.net/ch/auth/pkg/errors"
	"git.containerum.net/ch/kube-api/pkg/kubeErrors"
	"github.com/containerum/cherry"
	"github.com/containerum/chkit/pkg/model/configmap"
	"github.com/containerum/kube-client/pkg/cherry/resource-service"
	"github.com/sirupsen/logrus"
)

func (client *Client) CreateConfigMap(ns string, config configmap.ConfigMap) error {
	err := retry(4, func() (bool, error) {
		err := client.kubeAPIClient.CreateConfigMap(ns, config.Name, config.Data)
		switch {
		case err == nil:
			return false, nil
		case cherry.In(err,
			rserrors.ErrResourceNotExists()):
			return false, ErrResourceNotExists.CommentF("namespace %q doesn't exist", ns)
		case cherry.In(err,
			rserrors.ErrResourceNotOwned(),
			rserrors.ErrAccessRecordNotExists(),
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
		logrus.WithError(err).Errorf("unable to create configmap")
	}
	return err
}

func (client *Client) GetConfigmap(namespace, cmName string) (configmap.ConfigMap, error) {
	var gainedCM configmap.ConfigMap
	err := retry(4, func() (bool, error) {
		kubeConfigmap, err := client.kubeAPIClient.GetConfigMap(namespace, cmName)
		switch {
		case err == nil:
			gainedCM = configmap.ConfigMapFromKube(kubeConfigmap)
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
	return gainedCM, err
}

func (client *Client) GetConfigmapList(namespace string) (configmap.ConfigMapList, error) {
	var gainedCM configmap.ConfigMapList
	err := retry(4, func() (bool, error) {
		kubeConfigmapList, err := client.kubeAPIClient.GetConfigMapList(namespace)
		switch {
		case err == nil:
			gainedCM = configmap.ConfigMapListFromKube(kubeConfigmapList)
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
	return gainedCM, err
}
