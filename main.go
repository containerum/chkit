package main

import (
	"os"

	"github.com/containerum/chkit/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	defer angel()
	switch err := cmd.Run(os.Args).(type) {
	case nil:
		// pass
	default:
		logrus.Fatalf("Something bad happend: %v", err)
	}
}
