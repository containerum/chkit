package chClient

import (
	"git.containerum.net/ch/kube-client/pkg/cherry"
	"git.containerum.net/ch/kube-client/pkg/cherry/auth"
	"git.containerum.net/ch/kube-client/pkg/cherry/kube-api"
	"git.containerum.net/ch/kube-client/pkg/cherry/resource-service"
	"github.com/containerum/chkit/pkg/model/deployment"
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
			return false, ErrResourceNotExists
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
			rserrors.ErrResourceNotOwned(),
			rserrors.ErrAccessRecordNotExists()):
			return false, ErrYouDoNotHaveAccessToResource
		case cherry.In(err, autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
}
