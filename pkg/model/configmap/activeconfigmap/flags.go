package activeconfigmap

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/containerum/chkit/pkg/model/configmap"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/containerum/kube-client/pkg/model"
)

type Flags struct {
	Force      bool     `flag:"force f" desc:"suppress confirmation, optional"`
	File       string   `desc:"file with configmap data, .json, .yaml, .yml, optional"`
	Name       string   `desc:"configmap name, optional"`
	ItemFile   []string `flag:"item-file" desc:"configmap file, KEY:FILE_PATH or FILE_PATH"`
	ItemString []string `flag:"item-string" desc:"configmap item, KEY:VALUE string pair"`
}

func (flags Flags) ConfigMap() (configmap.ConfigMap, error) {
	var config = configmap.ConfigMap{
		Data: make(model.ConfigMapData, 16),
		Name: flags.Name,
	}

	if flags.File != "" {
		return FromFile(flags.File)
	} else if len(flags.ItemString) > 0 {
		items, err := getStringItems(flags.ItemString)
		if err != nil {
			return config, err
		}
		config = config.AddItems(items...)
	} else if len(flags.ItemFile) > 0 {
		items, err := getFileItems(flags.ItemFile)
		if err != nil {
			return config, err
		}
		config = config.AddItems(items...)
	}

	if config.Name == "" && len(config.Items()) > 0 {
		config.Name = namegen.Color() + "-" + config.Items()[0].Key()
	}

	return config, nil
}

func getStringItems(rawItems []string) ([]configmap.Item, error) {
	var items = make([]configmap.Item, 0, len(rawItems))
	for _, rawItem := range rawItems {
		var key string
		var value string
		if tokens := strings.SplitN(rawItem, ":", 2); len(tokens) == 2 {
			key = strings.TrimSpace(tokens[0])
			value = strings.TrimSpace(tokens[1])
		} else {
			return nil, fmt.Errorf("invalid token number in raw string item (got %v, required 2)", len(tokens))
		}
		items = append(items, configmap.NewItem(
			key,
			base64.StdEncoding.EncodeToString([]byte(value)),
		))
	}
	return items, nil
}

func getFileItems(rawItems []string) ([]configmap.Item, error) {
	var items = make([]configmap.Item, 0, len(rawItems))
	for _, rawItem := range rawItems {
		var filepath string
		var key string
		if tokens := strings.SplitN(rawItem, ":", 2); len(tokens) == 2 {
			key = strings.TrimSpace(tokens[0])
			filepath = tokens[1]
		} else if len(tokens) == 1 {
			key = path.Base(tokens[0])
			filepath = tokens[0]
		} else {
			return nil, fmt.Errorf("invalid token number in raw file item (got %v, required 2)", len(tokens))
		}
		value, err := ioutil.ReadFile(filepath)
		if err != nil {
			return nil, err
		}
		items = append(items, configmap.NewItem(
			key,
			string(value),
		))
	}
	return items, nil
}
