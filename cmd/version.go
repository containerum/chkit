package cmd

import (
	"runtime"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of chkit",
	Run: func(cmd *cobra.Command, args []string) {
		jww.FEEDBACK.Printf("Version: \nBuilt: \nOS: %s\nPlatform: %s\nGit commit: \n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
