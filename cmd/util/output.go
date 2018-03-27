package util

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/urfave/cli.v2"
)

func WriteData(ctx *cli.Context, renderer model.Renderer) error {
	var err error
	var data string
	switch ctx.String("output") {
	case "json":
		data, err = renderer.RenderJSON()
	case "yaml":
		data, err = renderer.RenderYAML()
	default:
		data = renderer.RenderTable()
	}
	if err != nil {
		return err
	}
	if !ctx.IsSet("file") {
		fmt.Println(data)
		return nil
	}
	return ioutil.WriteFile(ctx.String("file"), []byte(data), os.ModePerm)
}
