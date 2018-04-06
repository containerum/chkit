package util

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/urfave/cli.v2"
)

const (
	JSON = "json"
	YAML = "yaml"
)

func ExportDataCommand(ctx *cli.Context, renderer model.Renderer) error {
	var outputFile *string
	if ctx.IsSet("file") {
		f := ctx.String("file")
		outputFile = &f
	}
	return ExportData(ctx.String("output"), outputFile, renderer)
}

func ExportData(format string, outputFile *string, renderer model.Renderer) error {
	var data string
	var err error
	switch format {
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
	if outputFile == nil {
		fmt.Println(data)
		return nil
	}
	return ioutil.WriteFile(*outputFile, []byte(data), os.ModePerm)
}
