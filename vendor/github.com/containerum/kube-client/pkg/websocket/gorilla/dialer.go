// Package gorilla creates connection dialer for gorilla`s websocket implementation (github.com/gorilla/websocket)
package gorilla

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/containerum/cherry"
	"github.com/containerum/kube-client/pkg/identity"
	"github.com/gorilla/websocket"
)

type Dialer struct {
	token, fingerprint string
	serverAddr         *url.URL
	dialer             *websocket.Dialer
}

func WithToken(token string) func(ws *Dialer) {
	return func(ws *Dialer) {
		ws.token = token
	}
}

func WithFingerprint(fingerprint string) func(ws *Dialer) {
	return func(ws *Dialer) {
		ws.fingerprint = fingerprint
	}
}

func WithDialer(dialer *websocket.Dialer) func(ws *Dialer) {
	return func(ws *Dialer) {
		ws.dialer = dialer
	}
}

func WithHost(addr string) func(ws *Dialer) {
	return func(ws *Dialer) {
		var err error
		ws.serverAddr, err = url.Parse(addr)
		if err != nil {
			panic(err)
		}
		// check scheme
		switch ws.serverAddr.Scheme {
		case "ws", "wss":
			// pass
		case "http":
			ws.serverAddr.Scheme = "ws"
		case "https":
			ws.serverAddr.Scheme = "wss"
		default:
			panic(fmt.Errorf("invalid scheme \"%s\"", ws.serverAddr.Scheme))
		}
	}
}

func NewWebsocket(opts ...func(*Dialer)) *Dialer {
	ret := &Dialer{}

	for _, opt := range opts {
		opt(ret)
	}

	if ret.dialer == nil {
		ret.dialer = websocket.DefaultDialer
	}

	return ret
}

func (f *Dialer) SetToken(token string) {
	f.token = token
}

func (f *Dialer) SetFingerprint(fingerprint string) {
	f.fingerprint = fingerprint
}

func (f *Dialer) Dial(endpoint string, additionalHeaders http.Header) (*websocket.Conn, error) {
	reqURL, err := f.serverAddr.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set(identity.HeaderUserToken, f.token)
	headers.Set(identity.HeaderUserFingerprint, f.fingerprint)
	for key, value := range additionalHeaders {
		headers[key] = value
	}
	conn, resp, err := f.dialer.Dial(reqURL.String(), headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var chErr cherry.Err
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if unmarshalErr := json.Unmarshal(body, &chErr); unmarshalErr != nil {
			return nil, fmt.Errorf("unknown error: %s", body)
		}

		return nil, &chErr
	}

	return conn, nil
}
