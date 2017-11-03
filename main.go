package main

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/chlib"
	"github.com/containerum/chkit/cmd"
)

func main() {
	if _, err := os.Stat(chlib.ConfigDir); os.IsNotExist(err) {
		os.MkdirAll(chlib.ConfigDir, os.ModePerm)
		os.MkdirAll(chlib.TemplatesDir, os.ModePerm)
		os.MkdirAll(chlib.SrcFolder, os.ModePerm)
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
