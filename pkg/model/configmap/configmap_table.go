package configmap

import (
	"bytes"
	"fmt"

	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/util/text"
)

var (
	_ model.TableRenderer = ConfigMap{}
)

func (config ConfigMap) RenderTable() string {
	return model.RenderTable(config)
}

func (ConfigMap) TableHeaders() []string {
	return []string{
		"Name",
		"Age",
		"Items",
	}
}

func (config ConfigMap) TableRows() [][]string {
	return [][]string{{
		config.Name,
		config.Age(),
		func() string {
			buf := bytes.Buffer{}
			for _, item := range config.Items() {
				fmt.Fprint(&buf, text.Crop(item.String(), 32)+"\n")
			}
			return buf.String()
		}(),
	}}
}
