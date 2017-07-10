package cmd

import "github.com/spf13/cobra"

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure chkit default values",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	configCmd.PersistentFlags().StringP("set-token", "t", "", "Set user token")
	configCmd.PersistentFlags().StringP("set-default-namespace", "n", "default", "Default namespace")
	RootCmd.AddCommand(configCmd)
}
