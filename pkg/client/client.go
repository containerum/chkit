package chClient

import (
	kubeClient "git.containerum.net/ch/kube-client/pkg/client"
	"git.containerum.net/ch/kube-client/pkg/rest/re"
	"git.containerum.net/ch/kube-client/pkg/rest/remock"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"
)

const (
	// ErrUnableToInitClient -- unable to init client
	ErrUnableToInitClient chkitErrors.Err = "unable to init client"
)

// Client -- chkit core client
type Client struct {
	model.Config
	kubeAPIClient kubeClient.Client
}

// NewClient -- creates new client with provided options
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
		return nil, ErrUnableToInitClient.Wrap(err)
	}
	kubecli.SetFingerprint(config.Fingerprint)
	kubecli.SetToken(config.Tokens.AccessToken)
	chcli.kubeAPIClient = *kubecli
	for _, option := range options {
		chcli = option(chcli)
	}
	return chcli, nil
}

// UnsafeSkipTLSCheck -- optional client parameter to skip TLS verification
func UnsafeSkipTLSCheck(client *Client) *Client {
	restAPI := client.kubeAPIClient.RestAPI
	if _, ok := restAPI.(*re.Resty); ok || restAPI == nil {
		newRestAPI := re.NewResty(re.SkipTLSVerify)
		newRestAPI.SetFingerprint(client.Config.Fingerprint)
		newRestAPI.SetToken(client.Tokens.AccessToken)
		client.kubeAPIClient.RestAPI = newRestAPI
	}
	return client
}

// Mock -- optional parameter. Forces Client to use mock api
func Mock(client *Client) *Client {
	client.kubeAPIClient.RestAPI = remock.NewMock()
	return client
}
