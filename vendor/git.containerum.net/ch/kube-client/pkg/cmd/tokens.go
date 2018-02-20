package cmd

import (
	"fmt"
	"net/http"

	"git.containerum.net/ch/kube-client/pkg/model"
)

const (
	getCheckToken  = "/token/{access_token}"
	getExtendToken = "/token/{refresh_token}"
	userAgent      = "kube-client"
)

// CheckToken -- consumes JWT token, user fingerprint
// If they're correct returns user access data:
// list of namespaces and list of volumes OR uninitialized structure AND error
func (client *Client) CheckToken(token string) (model.CheckTokenResponse, error) {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"access_token": token,
		}).
		SetResult(model.CheckTokenResponse{}).
		Get(client.AuthURL + getCheckToken)
	if err != nil {
		return model.CheckTokenResponse{}, err
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		return *resp.Result().(*model.CheckTokenResponse), nil
	default:
		if resp.Error() != nil {
			return model.CheckTokenResponse{}, fmt.Errorf("%v", resp.Error())
		}
		return model.CheckTokenResponse{}, fmt.Errorf("%s", resp.Status())
	}
}

// ExtendToken -- consumes refresh JWT token and user fingerprint
// If they're correct returns new extended access and refresh token OR void tokens AND error.
// Old access and refresh token become inactive.
func (client *Client) ExtendToken(refreshToken string) (model.Tokens, error) {
	resp, err := client.Request.
		SetPathParams(map[string]string{
			"refresh_token": refreshToken,
		}).
		SetResult(model.Tokens{}).
		Put(client.AuthURL + getExtendToken)
	if err != nil {
		return model.Tokens{}, err
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		return *resp.Result().(*model.Tokens), nil
	default:
		if resp.Error() != nil {
			return model.Tokens{}, fmt.Errorf("%v", resp.Error())
		}
		return model.Tokens{}, fmt.Errorf("%s", resp.Status())
	}
}
