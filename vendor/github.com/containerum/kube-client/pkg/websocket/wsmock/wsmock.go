// Package wsmock implements mock websocket server
package wsmock

//go:generate noice -t errors.toml

import (
	"crypto/x509"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

type PeriodicServerConfig struct {
	MsgPeriod time.Duration
	MsgText   string
}

type PeriodicServer struct {
	cfg PeriodicServerConfig

	srv *httptest.Server
	log *logrus.Entry
}

func NewPeriodicServer(cfg PeriodicServerConfig, log *logrus.Entry, tls bool) *PeriodicServer {
	var ret PeriodicServer
	ret.cfg = cfg
	ret.log = log
	if tls {
		ret.srv = httptest.NewTLSServer(http.HandlerFunc(ret.periodicHandler))
	} else {
		ret.srv = httptest.NewServer(http.HandlerFunc(ret.periodicHandler))
	}
	return &ret
}

func (p *PeriodicServer) URL() *url.URL {
	serverURL, _ := url.Parse(p.srv.URL)
	switch serverURL.Scheme {
	case "http":
		serverURL.Scheme = "ws"
	case "https":
		serverURL.Scheme = "wss"
	}
	return serverURL
}

func (p *PeriodicServer) Certificate() *x509.Certificate {
	return p.srv.Certificate()
}

type EchoServer struct {
	srv *httptest.Server
	log *logrus.Entry
}

func NewEchoServer(log *logrus.Entry, tls bool) *EchoServer {
	var ret EchoServer
	ret.log = log
	if tls {
		ret.srv = httptest.NewTLSServer(http.HandlerFunc(ret.echoHandler))
	} else {
		ret.srv = httptest.NewServer(http.HandlerFunc(ret.echoHandler))
	}
	return &ret
}

func (s *EchoServer) URL() *url.URL {
	serverURL, _ := url.Parse(s.srv.URL)
	switch serverURL.Scheme {
	case "http":
		serverURL.Scheme = "ws"
	case "https":
		serverURL.Scheme = "wss"
	}

	return serverURL
}

func (s *EchoServer) Certificate() *x509.Certificate {
	return s.srv.Certificate()
}
