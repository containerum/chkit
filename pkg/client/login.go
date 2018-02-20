package client

import (
	"encoding/json"
	"io/ioutil"
	"os"

	kubeModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/err"
)

var (
	ErrUnableToLogin         = err.Err("unable to login")
	ErrUnableToPackTokens    = err.Err("unable to pack tokens")
	ErrUnableToSaveTokenFile = err.Err("unable to save token file")
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
	packedTokens, err := json.Marshal(tokens)
	if err != nil {
		return ErrUnableToPackTokens.Wrap(err)
	}
	err = ioutil.WriteFile(client.config.TokenFile,
		packedTokens,
		os.ModePerm)
	return err
}
