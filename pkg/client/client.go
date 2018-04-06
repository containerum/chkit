package chClient

import (
	kubeClient "git.containerum.net/ch/kube-client/pkg/client"
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
func NewClient(config model.Config, option KubeAPIclientFactory) (*Client, error) {
	chcli := &Client{
		Config: config,
	}
	var factory = WithCommonAPI
	if option != nil {
		factory = option
	}
	kcli, err := factory(config)
	if err != nil {
		return nil, ErrUnableToInitClient.Wrap(err)
	}
	chcli.kubeAPIClient = *kcli
	return chcli, nil
}
