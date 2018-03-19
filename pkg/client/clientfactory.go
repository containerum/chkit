package chClient

import (
	kubeClient "git.containerum.net/ch/kube-client/pkg/client"
	"git.containerum.net/ch/kube-client/pkg/rest/re"
	"git.containerum.net/ch/kube-client/pkg/rest/remock"
	"github.com/containerum/chkit/pkg/model"
)

type KubeAPIclientFactory func(model.Config) (*kubeClient.Client, error)

var (
	_ KubeAPIclientFactory = WithTestAPI
	_ KubeAPIclientFactory = WithMock
)

func WithCommonAPI(config model.Config) (*kubeClient.Client, error) {
	client, err := kubeClient.NewClient(kubeClient.Config{
		APIurl:  config.APIaddr,
		RestAPI: re.NewResty(),
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

func WithTestAPI(config model.Config) (*kubeClient.Client, error) {
	client, err := kubeClient.NewClient(kubeClient.Config{
		APIurl: config.APIaddr,
		User: kubeClient.User{
			Role: "user",
		},
	})
	if err != nil {
		return nil, err
	}
	newRestAPI := re.NewResty(re.SkipTLSVerify)
	newRestAPI.SetFingerprint(config.Fingerprint)
	newRestAPI.SetToken(config.Tokens.AccessToken)
	client.RestAPI = newRestAPI
	return client, nil
}

func WithMock(config model.Config) (*kubeClient.Client, error) {
	client, err := kubeClient.NewClient(kubeClient.Config{
		APIurl: config.APIaddr,
		User: kubeClient.User{
			Role: "user",
		},
	})
	if err != nil {
		return nil, err
	}
	newRestAPI := remock.NewMock()
	newRestAPI.SetFingerprint(config.Fingerprint)
	newRestAPI.SetToken(config.Tokens.AccessToken)
	client.RestAPI = newRestAPI
	return client, nil
}
