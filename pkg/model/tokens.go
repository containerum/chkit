package model

import "git.containerum.net/ch/kube-client/pkg/model"

var (
	_ = Tokens(model.Tokens{})
)

// Tokens -- access and refresh client tokens
type Tokens struct {
	AccessToken  string `json:"-"`
	RefreshToken string `json:"refresh_token"`
}
