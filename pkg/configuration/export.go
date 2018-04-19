package configuration

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/containerum/chkit/pkg/model"
)

type ExportFormat string

const (
	PRETTY ExportFormat = ""
	JSON   ExportFormat = "json"
	YAML   ExportFormat = "yaml"
)

type ExportConfig struct {
	Filename string
	Format   ExportFormat
}

func ExportData(renderer model.Renderer, config ExportConfig) error {
	var data string
	var err error
	switch config.Format {
	case JSON:
		data, err = renderer.RenderJSON()
	case YAML:
		data, err = renderer.RenderYAML()
	default:
		data = renderer.RenderTable()
	}
	if err != nil {
		return err
	}
	switch config.Filename {
	case "", "-":
		fmt.Println(data)
		return nil
	default:
		return ioutil.WriteFile(config.Filename, []byte(data), os.ModePerm)
	}
}
