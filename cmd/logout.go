package cmd

import (
	"github.com/containerum/chkit.v2/chlib"

	"github.com/spf13/cobra"
)

var logouCmd = &cobra.Command{
	Use:   "logout",
	Short: "Close session and remove token",
	Run: func(cmd *cobra.Command, args []string) {
		if err := chlib.UserLogout(db); err != nil {
			np.ERROR.Printf("Logout error: %s", err.Error())
		} else {
			np.FEEDBACK.Print("Bye!")
		}
	},
}

func init() {
	RootCmd.AddCommand(logouCmd)
}
