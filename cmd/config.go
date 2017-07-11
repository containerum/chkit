package cmd

import (
	"encoding/base64"
	"net/url"
	"time"

	"github.com/kfeofantov/chkit-v2/chlib"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure chkit default values",
	Run: func(cmd *cobra.Command, args []string) {
		info, err := chlib.GetUserInfo()
		if err != nil {
			jww.ERROR.Println(err)
			return
		}
		httpApi, err := chlib.GetHttpApiCfg()
		if err != nil {
			jww.ERROR.Println(err)
		}

		if cmd.Flags().NFlag() == 0 {
			jww.FEEDBACK.Println("Token: ", info.Token)
			jww.FEEDBACK.Println("Namespace: ", info.Namespace)
			jww.FEEDBACK.Println("HTTP API")
			jww.FEEDBACK.Println("\tServer: ", httpApi.Server)
			jww.FEEDBACK.Println("\tTimeout: ", httpApi.Timeout)
			return
		}

		if cmd.Flag("set-default-namespace").Changed {
			info.Namespace = cmd.Flag("set-default-namespace").Value.String()
			jww.FEEDBACK.Printf("Namespace changed to: %s\n", info.Namespace)
		}
		if cmd.Flag("set-token").Changed {
			enteredToken := cmd.Flag("set-token").Value.String()
			if _, err := base64.StdEncoding.DecodeString(enteredToken); err != nil {
				jww.ERROR.Println("Invalid token given")
				return
			}
			info.Token = enteredToken
			jww.FEEDBACK.Printf("Token changed to: %s\n", info.Token)
		}
		if cmd.Flag("set-http-server-address").Changed {
			address := cmd.Flag("set-http-server-address").Value.String()
			if _, err := url.ParseRequestURI(address); err != nil {
				jww.ERROR.Printf("Invalid http api server address given")
				return
			}
			httpApi.Server = address
		}
		if cmd.Flag("set-http-server-timeout").Changed {
			tm, err := cmd.Flags().GetDuration("set-http-server-timeout")
			if err != nil {
				jww.ERROR.Printf("Invalid http api timeout given")
				return
			}
			httpApi.Timeout = tm
		}

		err = chlib.UpdateUserInfo(info)
		if err != nil {
			jww.ERROR.Println(err)
		}
		err = chlib.UpdateHttpApiCfg(httpApi)
		if err != nil {
			jww.ERROR.Println(err)
		}
	},
}

func init() {
	configCmd.PersistentFlags().StringP("set-token", "t", "", "Set user token")
	configCmd.PersistentFlags().StringP("set-default-namespace", "n", "default", "Default namespace")
	configCmd.PersistentFlags().String("set-http-server-address", "http://146.185.135.181:3333", "HTTP API server address")
	configCmd.PersistentFlags().Duration("set-http-server-timeout", 10*time.Second, "HTTP API calls timeout")
	RootCmd.AddCommand(configCmd)
}
