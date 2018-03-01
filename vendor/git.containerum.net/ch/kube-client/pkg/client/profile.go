package client

import (
	"net/http"

	"git.containerum.net/ch/kube-client/pkg/cherry"

	"git.containerum.net/ch/kube-client/pkg/model"
)

const (
	userInfoPath           = "/user/info"
	userPasswordChangePath = "/password/change"
	userLoginPath          = "/login/basic"
)

// GetProfileInfo -- returns user info
func (client *Client) GetProfileInfo() (model.User, error) {
	resp, err := client.Request.
		SetResult(model.User{}).
		SetError(cherry.Err{}).
		Get(client.UserManagerURL + userInfoPath)
	if err := MapErrors(resp, err, http.StatusOK); err != nil {
		return model.User{}, err
	}
	return *resp.Result().(*model.User), nil
}

// ChangePassword -- changes user password, returns access and refresh tokens
func (client *Client) ChangePassword(currentPassword, newPassword string) (model.Tokens, error) {
	resp, err := client.Request.
		SetResult(model.Tokens{}).
		SetError(cherry.Err{}).
		Put(client.UserManagerURL + userPasswordChangePath)
	if err := MapErrors(resp, err, http.StatusAccepted, http.StatusOK); err != nil {
		return model.Tokens{}, err
	}
	return *resp.Error().(*model.Tokens), nil
}

// Login -- sign in with username and password
func (client *Client) Login(login model.Login) (model.Tokens, error) {
	resp, err := client.Request.
		SetBody(login).
		SetResult(model.Tokens{}).
		SetError(cherry.Err{}).
		Post(client.UserManagerURL + userLoginPath)
	if err = MapErrors(resp, err, http.StatusOK, http.StatusAccepted); err != nil {
		return model.Tokens{}, err
	}
	return *resp.Result().(*model.Tokens), nil

}
