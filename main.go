package main

import (
	"fmt"
	"os"

	"github.com/kfeofantov/chkit-v2/chlib"
	"github.com/kfeofantov/chkit-v2/cmd"
)

func main() {
	defer chlib.Close()
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
