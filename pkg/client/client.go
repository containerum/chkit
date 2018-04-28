package chClient

import (
	kubeClient "github.com/containerum/kube-client/pkg/client"
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
	return nil
}
