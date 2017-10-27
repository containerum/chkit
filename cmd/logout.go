package cmd

import (
	"github.com/spf13/cobra"
)

var logouCmd = &cobra.Command{
	Use:   "logout",
	Short: "Close session and remove token",
	Run: func(cmd *cobra.Command, args []string) {
		client.UserConfig.Token = ""
		saveUserSettings(*client.UserConfig)
		np.FEEDBACK.Print("Bye!")
	},
}

func init() {
	RootCmd.AddCommand(logouCmd)
}
