package chClient

import (
	"git.containerum.net/ch/auth/pkg/errors"
	"git.containerum.net/ch/kube-api/pkg/kubeErrors"
	permErrors "git.containerum.net/ch/permissions/pkg/errors"
	"git.containerum.net/ch/resource-service/pkg/rsErrors"
	"github.com/containerum/cherry"
	"github.com/containerum/chkit/pkg/model/configmap"
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

func (client *Client) DeleteConfigmap(namespace, cm string) error {
	var err = retry(4, func() (bool, error) {
		err := client.kubeAPIClient.DeleteConfigMap(namespace, cm)
		switch {
		case err == nil:
			return false, nil
		case cherry.In(err,
			rserrors.ErrResourceNotExists()):
			return false, ErrResourceNotExists.
				CommentF("service %q not found in %q", cm, namespace)
		case cherry.In(err,
			permErrors.ErrResourceNotOwned(),
			rserrors.ErrPermissionDenied()):

			return false, ErrYouDoNotHaveAccessToResource.
				CommentF("you don't have delete access to service %q", cm)
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			err = client.Auth()
			return true, err
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	if err != nil {
		logrus.WithError(err).Errorf("unable to delete configmap %q in %q", cm, namespace)
	}
	return err
}

func (client *Client) ReplaceConfigmap(namespaceID string, cm configmap.ConfigMap) error {
	var err = retry(4, func() (bool, error) {
		err := client.kubeAPIClient.UpdateConfigMap(namespaceID, cm.Name, cm.Copy().Data)
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
			err = client.Auth()
			return true, err
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	if err != nil {
		logrus.WithField("method", "ReplaceConfigmap").
			WithError(err).
			Errorf("unable to update configmap %q in %q", cm, namespaceID)
	}
	return err
}
