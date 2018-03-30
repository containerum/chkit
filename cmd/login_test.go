package cmd

import (
	"testing"

	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/model"
)

func TestLoginToSandbox(test *testing.T) {
	client, err := chClient.NewClient(model.Config{
		APIaddr: "https://192.168.88.200:8082",
		StorableConfig: model.StorableConfig{
			UserInfo: model.UserInfo{
				Username: "helpik94@yandex.ru",
				Password: "12345678",
			},
		},
		Fingerprint: Fingerprint(),
	}, chClient.WithTestAPI)
	if err != nil {
		test.Fatalf("error while client creation: %v", err)
	}
	if err := client.Login(); err != nil {
		test.Fatalf("error while login: %v", err)
	}
}
