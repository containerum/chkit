package chClient

import (
	"git.containerum.net/ch/auth/pkg/errors"
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
