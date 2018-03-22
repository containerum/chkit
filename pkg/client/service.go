package chClient

import (
	"git.containerum.net/ch/kube-client/pkg/cherry"
	"git.containerum.net/ch/kube-client/pkg/cherry/auth"
	"git.containerum.net/ch/kube-client/pkg/cherry/kube-api"
	"github.com/containerum/chkit/pkg/model/service"
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
