package chClient

import (
	kubeClient "git.containerum.net/ch/kube-client/pkg/client"
	kubeClientModels "git.containerum.net/ch/kube-client/pkg/model"
	"git.containerum.net/ch/kube-client/pkg/rest/re"
	"git.containerum.net/ch/kube-client/pkg/rest/remock"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"
)

type Client struct {
	Config        model.Config
	Tokens        kubeClientModels.Tokens
	kubeApiClient kubeClient.Client
}

func NewClient(config model.Config, options ...func(*Client) *Client) (*Client, error) {
	chcli := &Client{
		Config: config,
	}
	kubecli, err := kubeClient.NewClient(kubeClient.Config{
		APIurl:  config.APIaddr,
		RestAPI: re.NewResty(),
		User: kubeClient.User{
			Role: "user",
		},
	})
	if err != nil {
		err = chkitErrors.ErrUnableToInitClient().
			AddDetailsErr(err)
		return nil, err
	}
	kubecli.SetFingerprint(config.Fingerprint)
	chcli.kubeApiClient = *kubecli
	for _, option := range options {
		chcli = option(chcli)
	}
	return chcli, nil
}

func UnsafeSkipTLSCheck(client *Client) *Client {
	restAPI := client.kubeApiClient.RestAPI
	if _, ok := restAPI.(*re.Resty); ok || restAPI == nil {
		newRestAPI := re.NewResty(re.SkipTLSVerify)
		newRestAPI.SetFingerprint(client.Config.Fingerprint)
		client.kubeApiClient.RestAPI = newRestAPI
	}
	return client
}
func Mock(client *Client) *Client {
	client.kubeApiClient.RestAPI = remock.NewMock()
	return client
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
