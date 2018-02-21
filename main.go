package main

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/cmd"
)

func main() {
	if err := cmd.App.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if _, err := os.Stat(cmd.Configuration.ConfigPath); !os.IsNotExist(err) {
		return
	}
	if err := os.MkdirAll(cmd.Configuration.ConfigPath, os.ModePerm); err != nil {
		fmt.Printf("error while creating config dir: %v", err)
		os.Exit(1)
	}
}
