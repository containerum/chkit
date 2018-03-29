package main

import (
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v2"

	"github.com/containerum/chkit/cmd"
	"github.com/containerum/chkit/pkg/chkitErrors"
)

func main() {
	if !cmd.DEBUG {
		defer angel(recover())
	}
	switch err := cmd.Run(os.Args).(type) {
	case nil:
		// pass
	case chkitErrors.Err:
		fmt.Println(err)
	case cli.ExitCoder:
		fmt.Println(err)
	default:
		if !cmd.DEBUG {
			angel(err)
		}
	}
}
