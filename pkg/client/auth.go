package chClient

import (
	"git.containerum.net/ch/kube-client/pkg/cherry"
	"git.containerum.net/ch/kube-client/pkg/cherry/auth"
	"git.containerum.net/ch/kube-client/pkg/cherry/user-manager"
	kubeClientModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"
)

const (
	// ErrWrongPasswordLoginCombination -- wrong login-password combination
	ErrWrongPasswordLoginCombination chkitErrors.Err = "wrong login-password combination"
	// ErrUserNotExist -- user doesn't not exist
	ErrUserNotExist chkitErrors.Err = "user doesn't not exist"
)

// Auth -- refreshes tokens, on invalid token uses Login method to get new tokens
func (client *Client) Auth() error {
	if client.Tokens.RefreshToken != "" {
		err := client.Extend()
		switch {
		case err == nil:
			return nil
		case cherry.Equals(err, autherr.ErrInvalidToken()) ||
			cherry.Equals(err, autherr.ErrTokenNotFound()):
			break
		default:
			return err
		}
	}
	return client.Login()
}

// Login -- client login method. Updates tokens
func (client *Client) Login() error {
	tokens, err := client.kubeAPIClient.Login(kubeClientModels.Login{
		Login:    client.Config.Username,
		Password: client.Config.Password,
	})
	switch {
	case err == nil:
	case cherry.Equals(err, umErrors.ErrInvalidLogin()):
		return ErrWrongPasswordLoginCombination
	case cherry.Equals(err, umErrors.ErrUserNotExist()):
		return ErrUserNotExist
	default:
		return err
	}
	client.kubeAPIClient.SetToken(tokens.AccessToken)
	client.Tokens = model.Tokens(tokens)
	return nil
}

// Extend -- refreshes tokens, invalidates old
func (client *Client) Extend() error {
	tokens, err := client.kubeAPIClient.
		ExtendToken(client.Tokens.RefreshToken)
	if err != nil {
		return err
	}
	client.Tokens = model.Tokens(tokens)
	return nil
}
