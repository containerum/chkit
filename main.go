package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/containerum/chkit/cmd"
)

func main() {
	if err := cmd.Run(os.Args); err != nil {
		logrus.Fatalf("fatal error: %v", err)
	}
}
