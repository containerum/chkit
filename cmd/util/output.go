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
	output := "stdout"
	switch {
	case ctx.IsSet("json"):
		output = ctx.String("json")
		data, err = renderer.RenderJSON()
	case ctx.IsSet("yaml"):
		output = ctx.String("yaml")
		data, err = renderer.RenderYAML()
	default:
		data = renderer.RenderTable()
	}
	if err != nil {
		return err
	}
	if output == "stdout" {
		fmt.Println(data)
		return nil
	}
	return ioutil.WriteFile(output, []byte(data), os.ModePerm)
}
