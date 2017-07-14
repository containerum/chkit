package chlib

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadJsonFromFile(path string, b interface{}) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	err = json.NewDecoder(file).Decode(&b)
	return
}

func GetCmdRequestJson(client *Client, kind, name string) (ret []GenericJson, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("can`t extract field: %s", r)
		}
	}()
	var apiResult TcpApiResult
	switch kind {
	case KindNamespaces:
		apiResult, err = client.Get(KindNamespaces, name, "")
		if err != nil {
			return ret, err
		}
		items := apiResult["results"].([]interface{})
		for _, itemI := range items {
			item := itemI.(map[string]interface{})
			_, hasNs := item["data"].(map[string]interface{})["metadata"].(map[string]interface{})["namespace"]
			if hasNs {
				ret = append(ret, GenericJson(item))
			}
		}
	default:
		apiResult, err := client.Get(kind, name, client.userConfig.Namespace)
		if err != nil {
			return ret, err
		}
		items := apiResult["results"].([]interface{})
		for _, itemI := range items {
			ret = append(ret, itemI.(map[string]interface{}))
		}
	}
	return
}
