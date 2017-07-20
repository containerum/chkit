package main

import (
	"fmt"
	"os"

	"bytes"
	"chkit-v2/chlib"
	"chkit-v2/cmd"
	"runtime"
)

func main() {
	// remove developer`s paths from panic
	defer func() {
		if r := recover(); r != nil {
			println("Recovered error: ", r)
			// Get stack trace
			buf := make([]byte, 4096)
			buf = buf[:runtime.Stack(buf, false)]
			// Clean up stack trace
			buf = bytes.Replace(buf, []byte(chlib.DevGoPath), []byte{}, -1)
			buf = bytes.Replace(buf, []byte(chlib.DevGoRoot), []byte{}, -1)
			// Print stack trace
			println(string(buf))
		}
	}()
	if _, err := os.Stat(chlib.ConfigDir); os.IsNotExist(err) {
		os.MkdirAll(chlib.ConfigDir, os.ModePerm)
		os.MkdirAll(chlib.TemplatesFolder, os.ModePerm)
		os.MkdirAll(chlib.SrcFolder, os.ModePerm)
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
