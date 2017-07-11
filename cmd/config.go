package cmd

import (
	"encoding/base64"
	"fmt"

	"github.com/kfeofantov/chkit-v2/chlib"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure chkit default values",
	Run: func(cmd *cobra.Command, args []string) {
		info, err := chlib.GetUserInfo()
		if err != nil {
			fmt.Println(err)
			return
		}
		if cmd.Flags().NFlag() == 0 {
			fmt.Printf("Token: %s\nNamespace: %s\n", info.Token, info.Namespace)
			return
		}
		if cmd.Flag("set-default-namespace").Changed {
			info.Namespace = cmd.Flag("set-default-namespace").Value.String()
			fmt.Printf("Namespace changed to: %s\n", info.Namespace)
		}
		if cmd.Flag("set-token").Changed {
			enteredToken := cmd.Flag("set-token").Value.String()
			_, err := base64.StdEncoding.DecodeString(enteredToken)
			if err != nil {
				fmt.Println("Invalid token given")
				return
			}
			info.Token = enteredToken
			fmt.Printf("Token changed to: %s\n", info.Token)
		}
		err = chlib.UpdateUserInfo(info)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	configCmd.PersistentFlags().StringP("set-token", "t", "", "Set user token")
	configCmd.PersistentFlags().StringP("set-default-namespace", "n", "default", "Default namespace")
	RootCmd.AddCommand(configCmd)
}
