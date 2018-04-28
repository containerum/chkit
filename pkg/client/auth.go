package chClient

import (
	"git.containerum.net/ch/auth/pkg/errors"
	"github.com/containerum/kube-client/pkg/cherry/api-gateway"
	"git.containerum.net/ch/user-manager/pkg/umErrors"
	kubeClientModels "github.com/containerum/kube-client/pkg/model"
	"github.com/containerum/cherry"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"
	"github.com/sirupsen/logrus"
)

const (
	// ErrUnableToLogin -- unable to login
	ErrUnableToLogin chkitErrors.Err = "unable to login"
	// ErrUnableToRefreshToken -- unable to refresh token
	ErrUnableToRefreshToken chkitErrors.Err = "unable to refresh token"
	// ErrWrongPasswordLoginCombination -- wrong login-password combination
	ErrWrongPasswordLoginCombination chkitErrors.Err = "wrong login-password combination"
	// ErrUserNotExist -- user doesn't not exist
	ErrUserNotExist  chkitErrors.Err = "user doesn't not exist"
	ErrInternalError chkitErrors.Err = "internal server error"
)

// Auth -- refreshes tokens, on invalid token uses Login method to get new tokens
func (client *Client) Auth() error {
	if client.Tokens.RefreshToken != "" {
		logrus.Debugf("trying to extend token")
		err := client.Extend()
		switch {
		case err == nil:
			return nil
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound(),
			autherr.ErrTokenNotOwnedBySender()):
			logrus.WithError(err).Debugf("invalid token, trying to login")
			return client.Login()
		case cherry.In(err, gatewayErrors.ErrInternal()):
			logrus.WithError(ErrInternalError.Wrap(err)).
				Debugf("internal gateway error")
			return ErrInternalError
		default:
			logrus.WithError(ErrUnableToRefreshToken.Wrap(err)).
				Debugf("fatal auth error")
			return ErrUnableToRefreshToken
		}
	}
	return client.Login()
}

// Login -- client login method. Updates tokens
func (client *Client) Login() error {
	logrus.Debugf("start login")
	tokens, err := client.kubeAPIClient.Login(kubeClientModels.Login{
		Login:    client.Config.Username,
		Password: client.Config.Password,
	})
	switch {
	case err == nil:
	case cherry.Equals(err, umErrors.ErrInvalidLogin()):
		logrus.Debugf("invalid password login combination")
		return ErrWrongPasswordLoginCombination
	case cherry.Equals(err, umErrors.ErrUserNotExist()):
		logrus.Debugf("user does not exist")
		return ErrUserNotExist
	default:
		logrus.Debugf("fatal login error")
		return ErrUnableToLogin.Wrap(err)
	}
	client.kubeAPIClient.SetToken(tokens.AccessToken)
	client.Tokens = model.Tokens(tokens)
	return nil
}

// Extend -- refreshes tokens, invalidates old
func (client *Client) Extend() error {
	logrus.Debugf("extending tokens")
	tokens, err := client.kubeAPIClient.
		ExtendToken(client.Tokens.RefreshToken)
	if err != nil {
		logrus.Debugf("error while extending tokens")
		return err
	}
	client.Tokens = model.Tokens(tokens)
	client.kubeAPIClient.SetToken(tokens.AccessToken)
	return nil
}
