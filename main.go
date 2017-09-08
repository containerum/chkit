package main

import (
	"fmt"
	"os"

	"chkit-v2/chlib"
	"chkit-v2/cmd"
)

func main() {
	if _, err := os.Stat(chlib.ConfigDir); os.IsNotExist(err) {
		os.MkdirAll(chlib.ConfigDir, os.ModePerm)
		os.MkdirAll(chlib.TemplatesFolder, os.ModePerm)
		os.MkdirAll(chlib.SrcFolder, os.ModePerm)
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
