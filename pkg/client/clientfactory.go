package chClient

import (
	kubeClient "git.containerum.net/ch/kube-client/pkg/client"
	"git.containerum.net/ch/kube-client/pkg/rest/re"
	"git.containerum.net/ch/kube-client/pkg/rest/remock"
	"github.com/containerum/chkit/pkg/model"
)

// KubeAPIclientFactory -- creates new kube-client with provided config
type KubeAPIclientFactory func(model.Config) (*kubeClient.Client, error)

var (
	_ KubeAPIclientFactory = WithTestAPI
	_ KubeAPIclientFactory = WithMock
	_ KubeAPIclientFactory = WithCommonAPI
)

// WithCommonAPI -- creates kube-client for production api
func WithCommonAPI(config model.Config) (*kubeClient.Client, error) {
	client, err := kubeClient.NewClient(kubeClient.Config{
		APIurl:  config.APIaddr,
		RestAPI: re.NewResty(re.WithHost(config.APIaddr)),
		User: kubeClient.User{
			Role: "user",
		},
	})
	if err != nil {
		return nil, err
	}
	client.SetFingerprint(config.Fingerprint)
	client.SetToken(config.Tokens.AccessToken)
	return client, nil
}

// WithTestAPI -- creates kube-client for test api
func WithTestAPI(config model.Config) (*kubeClient.Client, error) {
	newRestAPI := re.NewResty(
		re.WithHost(config.APIaddr),
		re.SkipTLSVerify)
	newRestAPI.SetFingerprint(config.Fingerprint)
	newRestAPI.SetToken(config.Tokens.AccessToken)
	client, err := kubeClient.NewClient(kubeClient.Config{
		APIurl:  config.APIaddr,
		RestAPI: newRestAPI,
		User: kubeClient.User{
			Role: "user",
		},
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

// WithMock -- creates kube-client with mock API
func WithMock(config model.Config) (*kubeClient.Client, error) {
	newRestAPI := remock.NewMock()
	newRestAPI.SetFingerprint(config.Fingerprint)
	newRestAPI.SetToken(config.Tokens.AccessToken)
	client, err := kubeClient.NewClient(kubeClient.Config{
		APIurl:  config.APIaddr,
		RestAPI: newRestAPI,
		User: kubeClient.User{
			Role: "user",
		},
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
