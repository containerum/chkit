package chlib

import (
	"encoding/json"
	"fmt"
	"os"

	jww "github.com/spf13/jwalterweatherman"
	"gopkg.in/yaml.v2"
)

func LoadJsonFromFile(path string, b interface{}) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	err = json.NewDecoder(file).Decode(&b)
	return
}

func GetCmdRequestJson(client *Client, kind, name, nameSpace string) (ret []GenericJson, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("can`t extract field: %s", r)
		}
	}()
	apiResult, err := client.Get(kind, name, nameSpace)
	if err != nil {
		return ret, err
	}
	items := apiResult["results"].([]interface{})
	for _, itemI := range items {
		// remove kind "Namespace" results if namespace list requested
		if kind == KindNamespaces && name == "" &&
			itemI.(map[string]interface{})["data"].(map[string]interface{})["kind"].(string) == "Namespace" {
			continue
		}
		ret = append(ret, itemI.(map[string]interface{}))
	}
	return
}

func JsonPrettyPrint(jsonContent []GenericJson, np *jww.Notepad) (err error) {
	if len(jsonContent) == 0 {
		return fmt.Errorf("empty content received")
	}
	b, err := json.MarshalIndent(jsonContent[0]["data"], "", "    ")
	np.FEEDBACK.Printf("%s\n", b)
	return
}

func YamlPrint(jsonContent []GenericJson, np *jww.Notepad) (err error) {
	if len(jsonContent) == 0 {
		return fmt.Errorf("empty content received")
	}
	b, err := yaml.Marshal(jsonContent[0]["data"])
	np.FEEDBACK.Printf("%s\n", b)
	return
}
