package activeconfigmap

import (
	"fmt"
	"path/filepath"

	"encoding/json"
	"io/ioutil"

	"github.com/containerum/chkit/pkg/model/configmap"
	"gopkg.in/yaml.v2"
)

func FromFile(fname string) (configmap.ConfigMap, error) {
	var cm configmap.ConfigMap
	var content, err = ioutil.ReadFile(fname)
	if err != nil {
		return configmap.ConfigMap{}, err
	}
	switch filepath.Ext(fname) {
	case "yaml", "yml":
		return cm, yaml.Unmarshal(content, &cm)
	case "json":
		return cm, json.Unmarshal(content, &cm)
	default:
		return configmap.ConfigMap{}, fmt.Errorf("unknown file format of %q", fname)
	}
}
