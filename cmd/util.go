package cmd

import (
	"os"
)

func exitOnErr(err error) {
	if err != nil {
		log.WithError(err).Fatal(err)
		os.Exit(1)
	}
}
