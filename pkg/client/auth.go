package chClient

import (
	"git.containerum.net/ch/api-gateway/pkg/gatewayErrors"
	"git.containerum.net/ch/auth/pkg/errors"
	"git.containerum.net/ch/user-manager/pkg/umErrors"
	"github.com/containerum/cherry"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/util/coblog"
	kubeClientModels "github.com/containerum/kube-client/pkg/model"
)

const (
	// ErrUnableToLogin -- unable to login
	ErrUnableToLogin chkitErrors.Err = "unable to login"
	// ErrUnableToRefreshToken -- unable to refresh token
	// nolint:gas
	ErrUnableToRefreshToken chkitErrors.Err = "unable to refresh token"
	// ErrWrongPasswordLoginCombination -- wrong login-password combination
	ErrWrongPasswordLoginCombination chkitErrors.Err = "wrong login-password combination"
	// ErrUserNotExist -- user doesn't not exist
	ErrUserNotExist  chkitErrors.Err = "user doesn't not exist"
	ErrInternalError chkitErrors.Err = "internal server error"
)

// Auth -- refreshes tokens, on invalid token uses Login method to get new tokens
func (client *Client) Auth() error {
	var logger = coblog.Std.Component("client.Auth")
	logger.Debugf("START")
	defer logger.Debugf("END")
	if client.Tokens.RefreshToken != "" {
		logger.Debugf("trying to extend token")
		err := client.Extend()
		switch {
		case err == nil:
			logger.Debugf("OK")
			return nil
		case cherry.In(err,
			autherr.ErrInvalidToken(),
			autherr.ErrTokenNotFound(),
			autherr.ErrTokenNotOwnedBySender()):
			logger.WithError(err).Debugf("invalid token, trying to login")
			return client.Login()
		case cherry.In(err, gatewayErrors.ErrInternal()):
			logger.WithError(ErrInternalError.Wrap(err)).
				Debugf("internal gateway error")
			return ErrInternalError
		default:
			logger.WithError(ErrUnableToRefreshToken.Wrap(err)).
				Debugf("fatal auth error")
			return ErrUnableToRefreshToken
		}
	}
	logger.Debugf("empty refresh token, running login")
	return client.Login()
}

// Login -- client login method. Updates tokens
func (client *Client) Login() error {
	var logger = coblog.Std.Component("client.Login")
	logger.Debugf("START")
	defer logger.Debugf("END")
	logger.Debugf("running login as %q", client.Config.Username)
	tokens, err := client.kubeAPIClient.Login(kubeClientModels.Login{
		Login:    client.Config.Username,
		Password: client.Config.Password,
	})
	switch {
	case err == nil:
		logger.Debugf("OK")
		// pass
	case cherry.Equals(err, umErrors.ErrInvalidLogin()):
		logger.WithError(err).Errorf("invalid password login combination")
		return ErrWrongPasswordLoginCombination
	case cherry.Equals(err, umErrors.ErrUserNotExist()):
		logger.WithError(err).Errorf("user does not exist")
		return ErrUserNotExist
	default:
		logger.WithError(err).Errorf("fatal login error")
		return ErrUnableToLogin.Wrap(err)
	}
	logger.Debugf("setting tokens")
	client.kubeAPIClient.SetToken(tokens.AccessToken)
	client.Tokens = model.Tokens(tokens)
	return nil
}

// Extend -- refreshes tokens, invalidates old
func (client *Client) Extend() error {
	var logger = coblog.Std.Component("client.Extend tokens")
	logger.Debugf("START")
	defer logger.Debugf("END")
	logger.Debugf("extending tokens")
	tokens, err := client.kubeAPIClient.
		ExtendToken(client.Tokens.RefreshToken)
	if err != nil {
		logger.WithError(err).Errorf("unable to extend tokens")
		return err
	}
	logger.Debugf("setting tokens")
	client.Tokens = model.Tokens(tokens)
	client.kubeAPIClient.SetToken(tokens.AccessToken)
	return nil
}
