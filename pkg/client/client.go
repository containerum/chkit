package client

import (
	kubeClient "git.containerum.net/ch/kube-client/pkg/cmd"
	kubeModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/err"
	"github.com/containerum/chkit/pkg/model"
)

const (
	ErrUnableToCreateClient = err.Err("unable to create client")
)

type ChkitClient struct {
	kube   kubeClient.Client
	tokens kubeModels.Tokens
	config model.Config
}

func (client *ChkitClient) Config() model.Config {
	return client.config
}
func NewClient(config model.Config) (ChkitClient, error) {
	kube, err := kubeClient.CreateCmdClient(kubeClient.ClientConfig{
		APIurl:         config.APIaddr + ":1214",
		ResourceAddr:   config.APIaddr + ":1213",
		UserManagerURL: config.APIaddr + ":8111",
		User: kubeClient.User{
			Role: "user",
		},
	})
	if err != nil {
		return ChkitClient{}, ErrUnableToCreateClient.Wrap(err)
	}
	return ChkitClient{
		kube:   *kube,
		config: config,
	}, nil
}
