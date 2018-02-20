package client

import (
	kubeModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/err"
)

const (
	ErrUnableToLogin = err.Err("unable to login")
)

func (client *ChkitClient) Login(username, password string) error {
	tokens, err := client.kube.Login(kubeModels.Login{
		Username: username,
		Password: password,
	})
	if err != nil {
		return ErrUnableToLogin.Wrap(err)
	}
	client.tokens = tokens
	return nil
}
