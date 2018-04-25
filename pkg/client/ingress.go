package chClient

import (
	"git.containerum.net/ch/kube-client/pkg/cherry"
	"git.containerum.net/ch/kube-client/pkg/cherry/auth"
	"git.containerum.net/ch/kube-client/pkg/cherry/kube-api"
	"git.containerum.net/ch/kube-client/pkg/cherry/resource-service"
	"github.com/containerum/chkit/pkg/model/ingress"
	"github.com/sirupsen/logrus"
)

func (client *Client) GetIngress(ns, domain string) (ingress.Ingress, error) {
	var ingr ingress.Ingress
	err := retry(4, func() (bool, error) {
		kubeIngress, err := client.kubeAPIClient.GetIngress(ns, domain)
		switch {
		case err == nil:
			ingr = ingress.IngressFromKube(kubeIngress)
			return false, nil
		case cherry.In(err,
			kubeErrors.ErrResourceNotExist(),
			kubeErrors.ErrAccessError(),
			kubeErrors.ErrUnableGetResource()):
			return false, ErrResourceNotExists
		case cherry.In(err, autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	if err != nil {
		logrus.WithError(err).WithField("namespace", ns).
			Errorf("unable to get ingress")
	}
	return ingr, err
}

func (client *Client) GetIngressList(ns string) (ingress.IngressList, error) {
	var list ingress.IngressList
	err := retry(4, func() (bool, error) {
		kubeList, err := client.kubeAPIClient.GetIngressList(ns)
		switch {
		case err == nil:
			list = ingress.IngressListFromKube(kubeList)
			return false, nil
		case cherry.In(err,
			kubeErrors.ErrResourceNotExist(),
			kubeErrors.ErrAccessError(),
			kubeErrors.ErrUnableGetResource()):
			return false, ErrResourceNotExists
		case cherry.In(err, autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	if err != nil {
		logrus.WithError(err).WithField("namespace", ns).
			Errorf("unable to get ingress list")
	}
	return list, err
}

func (client *Client) AddIngress(ns string, ingr ingress.Ingress) error {
	return retry(4, func() (bool, error) {
		err := client.kubeAPIClient.AddIngress(ns, ingr.ToKube())
		switch {
		case err == nil:
			return false, nil
		case cherry.In(err,
			rserrors.ErrResourceNotExists()):
			logrus.WithError(ErrResourceNotExists.Wrap(err)).
				Debugf("error while creating ingress %q", ingr.Name)
			return false, ErrResourceNotExists.CommentF("namespace %q doesn't exist", ns)
		case cherry.In(err,
			rserrors.ErrResourceNotOwned(),
			rserrors.ErrAccessRecordNotExists(),
			rserrors.ErrPermissionDenied()):
			logrus.WithError(ErrYouDoNotHaveAccessToResource.Wrap(err)).
				Debugf("error while creating ingress %q", ingr.Name)
			return false, ErrYouDoNotHaveAccessToResource.
				CommentF("you don't have create access to namespace %q", ns)
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			err = client.Auth()
			if err != nil {
				logrus.WithError(err).
					Debugf("error while creating ingress %q", ingr.Name)
			}
			return true, err
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
}
