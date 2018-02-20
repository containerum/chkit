package client

import (
	"encoding/json"
	"io/ioutil"
	"os"

	kubeClient "git.containerum.net/ch/kube-client/pkg/cmd"
	kubeModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/err"
	"github.com/containerum/chkit/pkg/model"
)

const (
	ErrUnableToPackTokens    = err.Err("unable to pack tokens")
	ErrUnableToSaveTokenFile = err.Err("unable to save token file")
)

type ChkitClient struct {
	kube   kubeClient.Client
	tokens kubeModels.Tokens
	config model.Config
}

func (client *ChkitClient) SaveTokens() error {
	packedTokens, err := json.Marshal(client.tokens)
	if err != nil {
		return ErrUnableToPackTokens.Wrap(err)
	}
	err = ioutil.WriteFile(client.config.TokenFile,
		packedTokens,
		os.ModePerm)
	if err != nil {
		return ErrUnableToSaveTokenFile.Wrap(err)
	}
	return nil
}
