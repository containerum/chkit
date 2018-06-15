package chClient

import (
	"git.containerum.net/ch/auth/pkg/errors"
	"git.containerum.net/ch/kube-api/pkg/kubeErrors"
	permErrors "git.containerum.net/ch/permissions/pkg/errors"
	"git.containerum.net/ch/resource-service/pkg/rsErrors"
	"github.com/containerum/cherry"
	"github.com/containerum/chkit/pkg/model/deployment"
	kubeModels "github.com/containerum/kube-client/pkg/model"
	"github.com/sirupsen/logrus"
)

func (client *Client) GetDeployment(namespace, deplName string) (deployment.Deployment, error) {
	var depl deployment.Deployment
	err := retry(4, func() (bool, error) {
		kubeDeployment, err := client.kubeAPIClient.GetDeployment(namespace, deplName)
		switch {
		case err == nil:
			depl = deployment.DeploymentFromKube(kubeDeployment)
			return false, nil
		case cherry.In(err,
			kubeErrors.ErrResourceNotExist(),
			kubeErrors.ErrAccessError(),
			kubeErrors.ErrUnableGetResource()):
			return false, ErrResourceNotExists.Wrap(err)
		case cherry.In(err, autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	return depl, err
}

func (client *Client) GetDeploymentList(namespace string) (deployment.DeploymentList, error) {
	var list deployment.DeploymentList
	err := retry(4, func() (bool, error) {
		kubeList, err := client.kubeAPIClient.GetDeploymentList(namespace)
		switch {
		case err == nil:
			list = deployment.DeploymentListFromKube(kubeList)
			return false, nil
		case cherry.In(err,
			kubeErrors.ErrResourceNotExist(),
			kubeErrors.ErrAccessError(),
			kubeErrors.ErrUnableGetResource()):
			return false, ErrNamespaceNotExists
		case cherry.In(err, autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	return list, err
}

func (client *Client) DeleteDeployment(namespace, deplName string) error {
	return retry(4, func() (bool, error) {
		err := client.kubeAPIClient.DeleteDeployment(namespace, deplName)
		switch {
		case err == nil:
			return false, nil
		case cherry.In(err,
			rserrors.ErrResourceNotExists()):
			return false, ErrResourceNotExists.Wrap(err)
		case cherry.In(err,
			permErrors.ErrResourceNotOwned()):
			return false, ErrYouDoNotHaveAccessToResource
		case cherry.In(err, autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
}

func (client *Client) CreateDeployment(ns string, depl deployment.Deployment) error {
	return retry(4, func() (bool, error) {
		err := client.kubeAPIClient.CreateDeployment(ns, depl.ToKube())
		switch {
		case err == nil:
			return false, nil
		case cherry.In(err,
			rserrors.ErrResourceNotExists()):
			logrus.WithError(ErrResourceNotExists.Wrap(err)).
				Debugf("error while creating service %q", depl.Name)
			return false, ErrResourceNotExists.CommentF("namespace %q doesn't exist", ns)
		case cherry.In(err,
			permErrors.ErrResourceNotOwned(),
			rserrors.ErrPermissionDenied()):
			logrus.WithError(ErrYouDoNotHaveAccessToResource.Wrap(err)).
				Debugf("error while creating service %q", depl.Name)
			return false, ErrYouDoNotHaveAccessToResource.
				CommentF("you don't have create access to namespace %q", ns)
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			err = client.Auth()
			if err != nil {
				logrus.WithError(err).
					Debugf("error while creating deployment %q", depl.Name)
			}
			return true, err
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
}

func (client *Client) SetContainerImage(ns, depl string, image kubeModels.UpdateImage) error {
	return retry(4, func() (bool, error) {
		err := client.kubeAPIClient.SetContainerImage(ns, depl, image)
		switch {
		case err == nil:
			return false, nil
		case cherry.In(err,
			rserrors.ErrResourceNotExists()):
			logrus.WithError(ErrResourceNotExists.Wrap(err)).
				Errorf("unable to set image")
			return false, err
		case cherry.In(err,
			permErrors.ErrResourceNotOwned(),
			rserrors.ErrPermissionDenied()):
			logrus.WithError(ErrYouDoNotHaveAccessToResource.Wrap(err)).
				Errorf("unable to set container image")
			return false, ErrYouDoNotHaveAccessToResource.
				CommentF("you don't have create access to namespace %q", ns)
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			err = client.Auth()
			if err != nil {
				logrus.WithError(err).
					Errorf("unable to set container image")
			}
			return true, err
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
}

func (client *Client) ReplaceDeployment(ns string, newDepl deployment.Deployment) error {
	return retry(4, func() (bool, error) {
		err := client.kubeAPIClient.ReplaceDeployment(ns, newDepl.ToKube())
		switch {
		case err == nil:
			return false, nil
		case cherry.In(err,
			rserrors.ErrResourceNotExists()):
			logrus.WithError(ErrResourceNotExists.Wrap(err)).
				Debugf("error while creating service %q", newDepl.Name)
			return false, ErrResourceNotExists.Wrap(err)
		case cherry.In(err,
			permErrors.ErrResourceNotOwned(),
			rserrors.ErrPermissionDenied()):
			logrus.WithError(ErrYouDoNotHaveAccessToResource.Wrap(err)).
				Debugf("error while creating deployment %q", newDepl.Name)
			return false, ErrYouDoNotHaveAccessToResource.
				CommentF("you don't have create access to namespace %q", ns)
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			err = client.Auth()
			if err != nil {
				logrus.WithError(err).
					Debugf("error while creating service %q", newDepl.Name)
			}
			return true, err
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
}
