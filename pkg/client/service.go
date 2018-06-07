package chClient

import (
	"git.containerum.net/ch/auth/pkg/errors"
	"git.containerum.net/ch/kube-api/pkg/kubeErrors"
	permErrors "git.containerum.net/ch/permissions/pkg/errors"
	"git.containerum.net/ch/resource-service/pkg/rsErrors"
	"github.com/containerum/cherry"
	"github.com/containerum/chkit/pkg/model/service"
	"github.com/sirupsen/logrus"
)

func (client *Client) GetService(namespace, serviceName string) (service.Service, error) {
	var gainedService service.Service
	err := retry(4, func() (bool, error) {
		kubeService, err := client.kubeAPIClient.GetService(namespace, serviceName)
		switch {
		case err == nil:
			gainedService = service.ServiceFromKube(kubeService)
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
	return gainedService, err
}

func (client *Client) GetServiceList(namespace string) (service.ServiceList, error) {
	var gainedList service.ServiceList
	err := retry(4, func() (bool, error) {
		kubeLsit, err := client.kubeAPIClient.GetServiceList(namespace)
		switch {
		case err == nil:
			gainedList = service.ServiceListFromKube(kubeLsit)
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
	return gainedList, err
}

func (client *Client) DeleteService(namespace, service string) error {
	return retry(4, func() (bool, error) {
		err := client.kubeAPIClient.DeleteService(namespace, service)
		switch {
		case err == nil:
			return false, nil
		case cherry.In(err,
			rserrors.ErrResourceNotExists()):
			logrus.WithError(ErrResourceNotExists.Wrap(err)).
				Debugf("error while deleting service %q", service)
			return false, ErrResourceNotExists.
				CommentF("service %q not found in %q", service, namespace)
		case cherry.In(err,
			permErrors.ErrResourceNotOwned(),
			rserrors.ErrPermissionDenied()):
			logrus.WithError(ErrYouDoNotHaveAccessToResource.Wrap(err)).
				Debugf("error while deleting service %q", service)
			return false, ErrYouDoNotHaveAccessToResource.
				CommentF("you don't have delete access to service %q", service)
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			err = client.Auth()
			if err != nil {
				logrus.WithError(err).
					Debugf("error while deleting service %q", service)
			}
			return true, err
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
}

func (client *Client) CreateService(ns string, serv service.Service) error {
	return retry(4, func() (bool, error) {
		_, err := client.kubeAPIClient.CreateService(ns, serv.ToKube())
		switch {
		case err == nil:
			return false, nil
		case cherry.In(err,
			rserrors.ErrResourceNotExists()):
			logrus.WithError(ErrResourceNotExists.Wrap(err)).
				Debugf("error while creating service %q", serv.Name)
			return false, ErrResourceNotExists.Wrap(err)
		case cherry.In(err,
			permErrors.ErrResourceNotOwned(),
			rserrors.ErrPermissionDenied()):
			logrus.WithError(ErrYouDoNotHaveAccessToResource.Wrap(err)).
				Debugf("error while creating service %q", serv.Name)
			return false, ErrYouDoNotHaveAccessToResource.
				CommentF("you don't have create access to namespace %q", ns)
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			err = client.Auth()
			if err != nil {
				logrus.WithError(err).
					Debugf("error while creating service %q", serv.Name)
			}
			return true, err
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
}

func (client *Client) ReplaceService(ns string, serv service.Service) error {
	return retry(4, func() (bool, error) {
		_, err := client.kubeAPIClient.UpdateService(ns, serv.ToKube())
		switch {
		case err == nil:
			return false, nil
		case cherry.In(err,
			rserrors.ErrResourceNotExists()):
			logrus.WithError(ErrResourceNotExists.Wrap(err)).
				Debugf("error while replacing service %q", serv.Name)
			return false, ErrResourceNotExists.Wrap(err)
		case cherry.In(err,
			permErrors.ErrResourceNotOwned(),
			rserrors.ErrPermissionDenied()):
			logrus.WithError(ErrYouDoNotHaveAccessToResource.Wrap(err)).
				Debugf("error while replacing service %q", serv.Name)
			return false, ErrYouDoNotHaveAccessToResource.
				CommentF("you don't have write access to namespace %q", ns)
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			err = client.Auth()
			if err != nil {
				logrus.WithError(err).
					Debugf("error while replacing service %q", serv.Name)
			}
			return true, err
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
}
