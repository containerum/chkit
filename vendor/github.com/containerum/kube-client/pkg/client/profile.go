package client

import (
	"github.com/containerum/kube-client/pkg/rest"

	"github.com/containerum/kube-client/pkg/model"
)

const (
	userInfoPath           = "/user/info"
	userPasswordChangePath = "/password/change"
	userLoginPath          = "/login/basic"
)

// GetProfileInfo -- returns user info
func (client *Client) GetProfileInfo() (model.User, error) {
	var user model.User
	err := client.RestAPI.Get(rest.Rq{
		Result: &user,
		URL: rest.URL{
			Path:   userInfoPath,
			Params: rest.P{},
		},
	})
	return user, err
}

// ChangePassword -- changes user password, returns access and refresh tokens
func (client *Client) ChangePassword(currentPassword, newPassword string) (model.Tokens, error) {
	var tokens model.Tokens
	err := client.RestAPI.Put(rest.Rq{
		Result: &tokens,
		URL: rest.URL{
			Path:   userPasswordChangePath,
			Params: rest.P{},
		},
	})
	return tokens, err
}

// Login -- sign in with username and password
func (client *Client) Login(login model.Login) (model.Tokens, error) {
	var tokens model.Tokens
	err := client.RestAPI.Post(rest.Rq{
		Body:   login,
		Result: &tokens,
		URL: rest.URL{
			Path:   userLoginPath,
			Params: rest.P{},
		},
	})
	return tokens, err
}

// SetFingerprint -- sets fingerprint
func (client *Client) SetFingerprint(fingerprint string) {
	client.RestAPI.SetFingerprint(fingerprint)
	client.WSDialer.SetFingerprint(fingerprint)
}
