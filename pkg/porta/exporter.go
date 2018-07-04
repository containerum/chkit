package porta

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/containerum/chkit/pkg/model"
)

type Exporter struct {
	OutFile      string `desc:"output file" flag:"export-file"`
	OutputFormat string `desc:"output format, json/yaml" flag:"output o"`
}

func (exporter Exporter) ExporterActivated() bool {
	return exporter.OutFile != "" || exporter.OutputFormat != ""
}

func (exporter Exporter) renderFunc() func(renderer model.Renderer) ([]byte, error) {
	switch exporter.OutputFormat {
	case "json":
		return renderJSON
	case "yaml", "yml":
		return renderYAML
	default:
		var ext = path.Ext(exporter.OutFile)
		switch {
		case exporter.OutFile == "":
			return renderTable
		case ext == ".yaml", ext == ".yml":
			return renderYAML
		default:
			return renderJSON
		}
	}
}

func (exporter Exporter) Export(data model.Renderer) error {
	var b, err = exporter.renderFunc()(data)
	if err != nil {
		return err
	}
	if exporter.OutFile == "" {
		_, err = fmt.Printf("%s", b)
		return err
	}
	return ioutil.WriteFile(exporter.OutFile, b, os.ModePerm)
}

func renderJSON(renderer model.Renderer) ([]byte, error) {
	var data, err = renderer.RenderJSON()
	return []byte(data), err
}

func renderYAML(renderer model.Renderer) ([]byte, error) {
	var data, err = renderer.RenderYAML()
	return []byte(data), err
}

func renderTable(renderer model.Renderer) ([]byte, error) {
	var data = renderer.RenderTable()
	return []byte(data), nil
}
