package chClient

import (
	"crypto/tls"

	"time"

	"net"

	kubeClient "git.containerum.net/ch/kube-client/pkg/client"
	"git.containerum.net/ch/kube-client/pkg/rest/re"
	"git.containerum.net/ch/kube-client/pkg/rest/remock"
	"git.containerum.net/ch/kube-client/pkg/websocket/wsmock"
	"github.com/containerum/chkit/pkg/model"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
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
		WSDialer: websocket.DefaultDialer,
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
	var newRestAPI *re.Resty
	if config.Log == nil {
		newRestAPI = re.NewResty(
			re.WithHost(config.APIaddr),
			re.SkipTLSVerify)
	} else {
		newRestAPI = re.NewResty(
			re.WithLogger(config.Log),
			re.WithHost(config.APIaddr),
			re.SkipTLSVerify)
	}
	newRestAPI.SetFingerprint(config.Fingerprint)
	newRestAPI.SetToken(config.Tokens.AccessToken)
	client, err := kubeClient.NewClient(kubeClient.Config{
		APIurl:  config.APIaddr,
		RestAPI: newRestAPI,
		User: kubeClient.User{
			Role: "user",
		},
		WSDialer: &websocket.Dialer{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
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

	mockServer := wsmock.NewPeriodicServer(wsmock.PeriodicServerConfig{
		MsgPeriod: time.Second,
		MsgText:   "test\n",
	}, logrus.NewEntry(logrus.StandardLogger()), true)

	client, err := kubeClient.NewClient(kubeClient.Config{
		APIurl:  config.APIaddr,
		RestAPI: newRestAPI,
		User: kubeClient.User{
			Role: "user",
		},
		WSDialer: &websocket.Dialer{
			WriteBufferSize: 1024,
			ReadBufferSize:  1024,
			NetDial: func(network, addr string) (net.Conn, error) {
				return net.Dial("tcp", mockServer.URL().Host)
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
