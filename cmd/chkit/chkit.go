package main

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/pkg/cli"
)

func main() {
	if err := cli.Root(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
