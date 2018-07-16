package compose

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestModel(test *testing.T) {
	var compose = Compose{}
	var testDockerCompose, err = ioutil.ReadFile("docker-compose-test.yaml")
	if err != nil {
		test.Fatal(err)
	}
	if err := yaml.Unmarshal(testDockerCompose, &compose); err != nil {
		test.Fatal(err)
	}
	var data, _ = json.MarshalIndent(compose.Services["orderer.example.com"], "", "  ")
	test.Logf("\n%s", data)
}
