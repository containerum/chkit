package cmd

import (
	"fmt"
	"net/http"

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
		SetError(model.ResourceError{}).
		Get(client.UserManagerURL + userInfoPath)
	if err := catchErr(err, resp, http.StatusOK); err != nil {
		return model.User{}, err
	}
	return *resp.Result().(*model.User), nil
}

// ChangePassword -- changes user password, returns access and refresh tokens
func (client *Client) ChangePassword(currentPassword, newPassword string) (model.Tokens, error) {
	resp, err := client.Request.
		SetResult(model.Tokens{}).
		SetError(model.ResourceError{}).
		Put(client.UserManagerURL + userPasswordChangePath)
	if err := catchErr(err, resp, http.StatusAccepted, http.StatusOK); err != nil {
		return model.Tokens{}, err
	}
	return *resp.Error().(*model.Tokens), nil
}

// Login -- sign in with username and password
func (client *Client) Login(login model.Login) (model.Tokens, error) {
	resp, err := client.Request.
		SetBody(login).
		SetResult(model.Tokens{}).
		Post(client.UserManagerURL + userLoginPath)
	if err != nil {
		return model.Tokens{}, err
	}
	switch resp.StatusCode() {
	case http.StatusOK, http.StatusAccepted:
		return *resp.Result().(*model.Tokens), nil
	default:
		if resp.Error() != nil {
			return model.Tokens{}, fmt.Errorf("%v", resp.Error())
		}
		return model.Tokens{}, fmt.Errorf("%s", string(resp.Body()))
	}

}
