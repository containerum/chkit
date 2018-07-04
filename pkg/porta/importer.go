package porta

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type Importable interface {
	json.Unmarshaler
	yaml.Unmarshaler
}

type Importer struct {
	InFile      string `flag:"import-file"`
	InputFormat string `flag:"input i" desc:"input format, json/yaml"`
}

func (importer Importer) ImportActivated() bool {
	return importer.InputFormat != "" || importer.InFile != ""
}

func (importer Importer) Import(obj Importable) error {
	var data []byte
	if importer.InFile == "" {
		var buf = &bytes.Buffer{}
		if _, err := buf.ReadFrom(os.Stdin); err != nil {
			return err
		}
		data = buf.Bytes()
	} else {
		var err error
		data, err = ioutil.ReadFile(importer.InFile)
		if err != nil {
			return err
		}
	}
	return importer.importFunc()(data, obj)
}

func (importer Importer) importFunc() func(data []byte, obj Importable) error {
	switch importer.InputFormat {
	case "json":
		return importJSON
	case "yaml", "yml":
		return importYAML
	default:
		var ext = path.Ext(importer.InFile)
		switch {
		case ext == ".yaml", ext == ".yml":
			return importYAML
		default:
			return importJSON
		}
	}
}

func importJSON(data []byte, obj Importable) error {
	return obj.UnmarshalJSON(data)
}

func importYAML(data []byte, obj Importable) error {
	return yaml.Unmarshal(data, obj)
}
