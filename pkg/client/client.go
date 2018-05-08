package chClient

import (
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model"
	kubeClient "github.com/containerum/kube-client/pkg/client"
)

const (
	// ErrUnableToInitClient -- unable to init client
	ErrUnableToInitClient chkitErrors.Err = "unable to init client"
)

// Client -- chkit core client
type Client struct {
	model.Config
	factory       KubeAPIclientSetup
	isInitialized bool
	kubeAPIClient kubeClient.Client
}

func (client *Client) IsInitialized() bool {
	return client.isInitialized
}

func (client *Client) Init(setup KubeAPIclientSetup) error {
	if setup == nil {
		setup = WithCommonAPI
	}
	kcli, err := setup(client.Config)
	if err != nil {
		return ErrUnableToInitClient.Wrap(err)
	}
	client.kubeAPIClient = *kcli
	client.isInitialized = true
	client.factory = setup
	return nil
}

func (client *Client) ReInit() error {
	if client.factory == nil || !client.isInitialized {
		panic("[chkit/pkg/client.Client.ReInit] try to reinit not initialized client")
	}
	client.Init(client.factory)
	return nil
}
