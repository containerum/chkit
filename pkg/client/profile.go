package chClient

import (
	"github.com/containerum/chkit/pkg/model/user"
	"github.com/sirupsen/logrus"
)

func (client *Client) GetProfile() (user.User, error) {
	var usr user.User
	err := retry(4, func() (bool, error) {
		kubeUsr, err := client.kubeAPIClient.GetProfileInfo()
		if err == nil {
			usr = user.UserFromKube(kubeUsr)
		}
		return HandleErrorRetry(client, err)
	})
	if err != nil {
		logrus.WithError(err).
			Errorf("unable to get profile info")
	}
	return usr, err
}
