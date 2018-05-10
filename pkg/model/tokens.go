package model

import "github.com/containerum/kube-client/pkg/model"

var (
	_ = Tokens(model.Tokens{})
)

// Tokens -- access and refresh client tokens
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
