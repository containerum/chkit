package cmd

import (
	"net/url"
	"os"

	"github.com/go-resty/resty"
)

//TODO: Make Interface

//Client - rest client
type Client struct {
	*resty.Request
	ClientConfig
	User User
}

//User -
type User struct {
	Role string
}

// ClientConfig -- provides configuration for Client
// If APIurl or ResourceAddr is void,
// trys to get them from envvars
type ClientConfig struct {
	User           User
	APIurl         string
	ResourceAddr   string
	UserManagerURL string
	AuthURL        string
}

//CreateCmdClient -
func CreateCmdClient(config ClientConfig) (*Client, error) {
	var APIurl *url.URL
	var err error
	if config.APIurl == "" {
		APIurl, err = url.Parse(os.Getenv("API_URL"))
	} else {
		APIurl, err = url.Parse(config.APIurl)
	}
	if err != nil {
		return nil, err
	}
	config.APIurl = APIurl.String()

	if config.ResourceAddr == "" {
		// TODO: addr validation
		config.ResourceAddr = os.Getenv("RESOURCE_ADDR")
	}
	if config.UserManagerURL == "" {
		config.UserManagerURL = os.Getenv("USER_MANAGER_URL")
	}
	if config.AuthURL == "" {
		config.AuthURL = os.Getenv("AUTH_URL")
	}
	client := &Client{
		Request:      resty.R(),
		ClientConfig: config,
		User:         config.User,
	}
	client.SetHeaders(map[string]string{
		"X-User-Role": client.User.Role,
	})
	return client, nil
}
