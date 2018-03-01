package client

import (
	kubeClient "git.containerum.net/ch/kube-client/pkg/client"
	kubeClientModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"
)

type Client struct {
	tokens        kubeClientModels.Tokens
	kubeApiClient kubeClient.Client
}

func NewClient(config model.Config) (*Client, error) {
	chcli := &Client{}
	kubecli, err := kubeClient.CreateCmdClient(kubeClient.Config{
		APIurl:         config.APIaddr + ":1214",
		UserManagerURL: config.APIaddr + ":8111",
		ResourceAddr:   config.APIaddr + ":1213",
		AuthURL:        config.APIaddr + ":1111",
		User: kubeClient.User{
			Role: "user",
		},
	})
	if err != nil {
		err = chkitErrors.ErrUnableToInitClient().
			AddDetailsErr(err)
		return nil, err
	}
	chcli.kubeApiClient = *kubecli
	return chcli, nil
}

func (client *Client) Login(username, password string) error {
	tokens, err := client.kubeApiClient.Login(kubeClientModels.Login{
		Username: username,
		Password: password,
	})
	if err != nil {
		return err
	}
	client.tokens = tokens
	return nil
}
