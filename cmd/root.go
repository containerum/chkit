package cmd

import (
	"os"
)

func exitOnErr(err error) {
	if err != nil {
		notepad.ERROR.Println(err)
		os.Exit(1)
	}
}
