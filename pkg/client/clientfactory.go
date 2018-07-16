package chClient

import (
	"crypto/tls"

	"time"

	"github.com/containerum/chkit/pkg/model"
	kubeClient "github.com/containerum/kube-client/pkg/client"
	"github.com/containerum/kube-client/pkg/rest/re"
	"github.com/containerum/kube-client/pkg/rest/remock"
	"github.com/containerum/kube-client/pkg/websocket/gorilla"
	"github.com/containerum/kube-client/pkg/websocket/wsmock"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// KubeAPIclientSetup -- creates new kube-client with provided config
type KubeAPIclientSetup func(model.Config) (*kubeClient.Client, error)

var (
	_ KubeAPIclientSetup = WithTestAPI
	_ KubeAPIclientSetup = WithMock
	_ KubeAPIclientSetup = WithCommonAPI
)

// WithCommonAPI -- creates kube-client for production api
func WithCommonAPI(config model.Config) (*kubeClient.Client, error) {
	client, err := kubeClient.NewClient(kubeClient.Config{
		RestAPI: re.NewResty(re.WithHost(config.APIaddr)),
		User: kubeClient.User{
			Role: "user",
		},
		WSDialer: gorilla.NewWebsocket(gorilla.WithHost(config.APIaddr)),
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
	dialer := &websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// nolint:gas
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client, err := kubeClient.NewClient(kubeClient.Config{
		// nolint:gas
		RestAPI: re.NewResty(re.WithHost(config.APIaddr), re.WithLogger(config.Log), re.SkipTLSVerify),
		User: kubeClient.User{
			Role: "user",
		},
		WSDialer: gorilla.NewWebsocket(gorilla.WithHost(config.APIaddr), gorilla.WithDialer(dialer)),
	})
	if err != nil {
		return nil, err
	}
	client.SetFingerprint(config.Fingerprint)
	client.SetToken(config.Tokens.AccessToken)
	return client, nil
}

// WithMock -- creates kube-client with mock API
func WithMock(config model.Config) (*kubeClient.Client, error) {
	mockServer := wsmock.NewPeriodicServer(wsmock.PeriodicServerConfig{
		MsgPeriod: time.Second,
		MsgText:   "test\n",
	}, logrus.NewEntry(logrus.StandardLogger()), true)

	dialer := &websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// nolint:gas
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client, err := kubeClient.NewClient(kubeClient.Config{
		RestAPI: remock.NewMock(),
		User: kubeClient.User{
			Role: "user",
		},
		WSDialer: gorilla.NewWebsocket(gorilla.WithHost(mockServer.URL().String()), gorilla.WithDialer(dialer)),
	})
	if err != nil {
		return nil, err
	}
	client.SetToken(config.Tokens.AccessToken)
	client.SetFingerprint(config.Fingerprint)
	return client, nil
}
