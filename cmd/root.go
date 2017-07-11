package cmd

import (
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

//RootCmd main cmd entrypoint
var RootCmd = &cobra.Command{
	Use: "chkit",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			cmd.Usage()
		}
		var err error
		debug, err = cmd.Flags().GetBool("debug")
		if err != nil {
			jww.ERROR.Println(err)
			return
		}
	},
}

var debug bool

func init() {
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "turn on debugging messages")
}
