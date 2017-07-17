package cmd

import "github.com/spf13/cobra"

//RootCmd main cmd entrypoint
var RootCmd = &cobra.Command{
	Use: "chkit",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			cmd.Usage()
		}
		debug, _ = cmd.Flags().GetBool("debug")
	},
}

var debug bool

func init() {
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "turn on debugging messages")
}
