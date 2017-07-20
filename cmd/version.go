package cmd

import (
	"runtime"
	"time"

	"chkit-v2/chlib"
	"chkit-v2/helpers"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of chkit",
	Run: func(cmd *cobra.Command, args []string) {
		np.FEEDBACK.Printf("CH Client\n Version: %s\n Built: %s\n OS: %s\n Platform: %s\n Git commit: %s\n",
			helpers.CurrentClientVersion, chlib.BuildDate, runtime.GOOS, runtime.GOARCH, chlib.CommitHash)
	},
}

func init() {
	checkBuildDate()
	RootCmd.AddCommand(versionCmd)
}

func checkBuildDate() {
	if chlib.BuildDate == "" {
		t := helpers.GetProgramBuildTime()
		chlib.BuildDate = t.Format(time.RFC822)
		return
	}
	if t, err := time.Parse("2006-01-02T15:04:05MST", chlib.BuildDate); err == nil {
		chlib.BuildDate = t.Format(time.RFC822)
	}
}
