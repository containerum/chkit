package main

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/pkg/cli"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
)

func main() {
	defer func() {
		switch panicInfo := recover().(type) {
		case nil:
			// pass
		default:
			angel.Angel(&context.Context{Version: cli.VERSION}, panicInfo)
		}
	}()
	if err := cli.Root(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
