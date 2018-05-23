package chClient

import (
	"git.containerum.net/ch/auth/pkg/errors"
	"git.containerum.net/ch/kube-api/pkg/kubeErrors"
	"git.containerum.net/ch/user-manager/pkg/umErrors"
	"github.com/containerum/cherry"
	"github.com/containerum/chkit/pkg/model/user"
	"github.com/sirupsen/logrus"
)

func (client *Client) GetProfile() (user.User, error) {
	var usr user.User
	err := retry(4, func() (bool, error) {
		kubeUsr, err := client.kubeAPIClient.GetProfileInfo()
		switch {
		case err == nil:
			usr = user.UserFromKube(kubeUsr)
			return false, nil
		case cherry.In(err,
			umErrors.ErrAccountBlocked(),
			umErrors.ErrPermissionsError(),
			umErrors.ErrNotActivated(),
			kubeErrors.ErrUnableGetResource()):
			return false, err
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound(),
			autherr.ErrTokenNotOwnedBySender()):
			return true, client.Auth()
		default:
			return true, ErrFatalError.Wrap(err)
		}
	})
	if err != nil {
		logrus.WithError(err).
			Errorf("unable to get profile info")
	}
	return usr, err
}
