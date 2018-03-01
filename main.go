package main

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/cmd"
)

func main() {
	if err := cmd.App.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
