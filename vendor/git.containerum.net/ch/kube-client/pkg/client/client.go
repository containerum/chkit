package client

import (
	"net/url"

	"git.containerum.net/ch/kube-client/pkg/rest"
)

//TODO: Make Interface

//Client - rest client
type Client struct {
	Config
}

//User -
type User struct {
	Role string
}

// Config -- provides configuration for Client
// If APIurl or ResourceAddr is void,
// trys to get them from envvars
type Config struct {
	User    User
	APIurl  string
	RestAPI rest.REST
}

//NewClient -
func NewClient(config Config) (*Client, error) {
	var APIurl *url.URL
	var err error
	APIurl, err = url.Parse(config.APIurl)
	if err != nil {
		return nil, err
	}
	config.APIurl = APIurl.String()
	if config.RestAPI == nil {
		panic("[kube-client] undefined RestAPI in config")
	}
	client := &Client{
		config,
	}
	return client, nil
}
