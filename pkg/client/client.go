package chClient

import (
	"git.containerum.net/ch/kube-client/pkg/cherry"
	kubeClient "git.containerum.net/ch/kube-client/pkg/client"
	kubeClientModels "git.containerum.net/ch/kube-client/pkg/model"
	"git.containerum.net/ch/kube-client/pkg/rest/re"
	"git.containerum.net/ch/kube-client/pkg/rest/remock"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"
)

const (
	// ErrUnableToInitClient -- unable to init client
	ErrUnableToInitClient chkitErrors.Err = "unable to init client"
	// ErrWrongPasswordLoginCombination -- wrong login-password combination
	ErrWrongPasswordLoginCombination chkitErrors.Err = "wrong login-password combination"
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

func (client *Client) Auth() error {
	if client.Tokens.RefreshToken != "" {
		switch err := client.Extend().(type) {
		case *cherry.Err:
			switch err.ID.Kind {
			case 2, 3:
				// if token is rotten, then login
			default:
				return err
			}
		default:
			return err
		}

	}
	switch err := client.Login().(type) {
	case *cherry.Err:
		switch err.ID.Kind {
		case 6, 19:
			return ErrWrongPasswordLoginCombination
		default:
			return err
		}
	default:
		return err
	}
}

// Login -- client login method. Updates tokens
func (client *Client) Login() error {
	tokens, err := client.kubeAPIClient.Login(kubeClientModels.Login{
		Login:    client.Config.Username,
		Password: client.Config.Password,
	})
	if err != nil {
		return err
	}
	client.kubeAPIClient.SetToken(tokens.AccessToken)
	client.Tokens = model.Tokens(tokens)
	return nil
}

func (client *Client) Extend() error {
	tokens, err := client.kubeAPIClient.
		ExtendToken(client.Tokens.RefreshToken)
	if err != nil {
		return err
	}
	client.Tokens = model.Tokens(tokens)
	return nil
}
