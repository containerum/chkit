package main

import (
	"net"
	"net/url"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/containerum/chkit/cmd"
)

func main() {
	switch err := cmd.Run(os.Args).(type) {
	case nil:
	case *url.Error:
		logrus.Printf("%T", err.Err)
		switch err := err.Err.(type) {
		case *net.OpError:
			logrus.Printf("%T", err.Err)
		}
	default:
		logrus.Fatalf("Something bad happend: %v", err)
	}
}
