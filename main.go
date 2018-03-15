package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/containerum/chkit/cmd"
)

func main() {
	switch err := cmd.Run(os.Args).(type) {
	case nil:
	default:
		logrus.Fatalf("Something bad happend: %v", err)
	}
}
