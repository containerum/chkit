package cmd

import (
	"github.com/kfeofantov/chkit-v2/chlib"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var logouCmd = &cobra.Command{
	Use:   "logout",
	Short: "Close session and remove token",
	Run: func(cmd *cobra.Command, args []string) {
		if err := chlib.UserLogout(); err != nil {
			jww.ERROR.Printf("Logout error: %s", err.Error())
		} else {
			jww.FEEDBACK.Print("Bye!")
		}
	},
}

func init() {
	RootCmd.AddCommand(logouCmd)
}
