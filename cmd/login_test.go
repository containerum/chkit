package cmd

import (
	"testing"

	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/model"
)

func TestLoginToSandbox(test *testing.T) {
	client, err := chClient.NewClient(model.ClientConfig{
		APIaddr:  "https://192.168.88.200:8082",
		Username: "helpik94@yandex.ru",
		Password: "12345678",
	}, chClient.UnsafeSkipTLSCheck)
	if err != nil {
		test.Fatalf("error while client creation: %v", err)
	}
	if err := client.Login(); err != nil {
		test.Fatalf("error while login: %v", err)
	}
}
