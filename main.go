package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/containerum/chkit/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		recoverData := recover()
		switch recoverData := recoverData.(type) {
		case nil:
		case error:
			fmt.Printf("[FATAL] %v", recoverData)
		default:
			fmt.Printf("[FATAL] %v\n%s", recoverData, string(debug.Stack()))
		}
	}()
	switch err := cmd.Run(os.Args).(type) {
	case nil:
	default:
		logrus.Fatalf("Something bad happend: %v", err)
	}
}
