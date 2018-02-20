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
	serverURL string
	User      User
}

//User -
type User struct {
	Role string
}

//CreateCmdClient -
func CreateCmdClient(u User) (*Client, error) {
	apiURL, err := url.Parse(os.Getenv("API_URL"))
	if err != nil {
		return nil, err
	}
	client := &Client{
		Request:   resty.R(),
		serverURL: apiURL.String(),
		User:      u,
	}
	client.SetHeaders(map[string]string{
		"X-User-Role": client.User.Role,
	})
	return client, nil
}
