package chClient

import (
	kubeClient "git.containerum.net/ch/kube-client/pkg/client"
	kubeClientModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"
)

type Client struct {
	Config        model.ClientConfig
	Tokens        kubeClientModels.Tokens
	kubeApiClient kubeClient.Client
}

func NewClient(config model.ClientConfig) (*Client, error) {
	chcli := &Client{
		Config: config,
	}
	kubecli, err := kubeClient.NewClient(kubeClient.Config{
		APIurl: config.APIaddr,
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

func (client *Client) Login() error {
	tokens, err := client.kubeApiClient.Login(kubeClientModels.Login{
		Username: client.Config.Username,
		Password: client.Config.Password,
	})
	if err != nil {
		return err
	}
	client.kubeApiClient.SetToken(tokens.AccessToken)
	client.Tokens = tokens
	return nil
}
