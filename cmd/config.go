package cmd

import (
	"encoding/base64"
	"net"
	"net/url"

	"chkit-v2/chlib"
	"chkit-v2/chlib/dbconfig"
	"chkit-v2/helpers"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"os"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure chkit default values",
	Run: func(cmd *cobra.Command, args []string) {
		info, err := db.GetUserInfo()
		if err != nil {
			jww.ERROR.Println(err)
			return
		}
		httpApi, err := db.GetHttpApiConfig()
		if err != nil {
			jww.ERROR.Println(err)
			return
		}
		tcpApi, err := db.GetTcpApiConfig()
		if err != nil {
			jww.ERROR.Println(err)
			return
		}

		if cmd.Flags().NFlag() == 0 {
			jww.FEEDBACK.Println("Token: ", info.Token)
			jww.FEEDBACK.Println("Namespace: ", info.Namespace)
			jww.FEEDBACK.Println("HTTP API")
			jww.FEEDBACK.Println("\tServer: ", httpApi.Server)
			jww.FEEDBACK.Println("\tTimeout: ", httpApi.Timeout)
			jww.FEEDBACK.Println("TCP API")
			jww.FEEDBACK.Printf("\tServer: %s", tcpApi.Address)
			jww.FEEDBACK.Println("\tBuffer size: ", tcpApi.BufferSize)
			return
		}

		if cmd.Flag("set-default-namespace").Changed {
			newNamespace := cmd.Flag("set-default-namespace").Value.String()
			client, err := chlib.NewClient(db, helpers.CurrentClientVersion, helpers.UuidV4())
			if err != nil {
				jww.ERROR.Println(err)
				os.Exit(1)
			}
			if _, err := client.Get(chlib.KindNamespaces, newNamespace, ""); err != nil {
				jww.ERROR.Println(err)
				os.Exit(1)
			}
			jww.FEEDBACK.Printf("Namespace changed to: %s\n", info.Namespace)
		}
		if cmd.Flag("set-token").Changed {
			enteredToken := cmd.Flag("set-token").Value.String()
			if _, err := base64.StdEncoding.DecodeString(enteredToken); err != nil {
				jww.FEEDBACK.Println("Invalid token given")
				os.Exit(1)
			}
			info.Token = enteredToken
			jww.FEEDBACK.Printf("Token changed to: %s\n", info.Token)
		}
		if cmd.Flag("set-http-server-address").Changed {
			address := cmd.Flag("set-http-server-address").Value.String()
			if _, err := url.ParseRequestURI(address); err != nil {
				jww.FEEDBACK.Printf("Invalid HTTP API server address given")
				os.Exit(1)
			}
			httpApi.Server = address
			jww.FEEDBACK.Printf("HTTP API server address changed to: %s", address)
		}
		if cmd.Flag("set-http-server-timeout").Changed {
			tm, err := cmd.Flags().GetDuration("set-http-server-timeout")
			if err != nil {
				jww.FEEDBACK.Printf("Invalid HTTP API timeout given")
				os.Exit(1)
			}
			httpApi.Timeout = tm
			jww.FEEDBACK.Printf("HTTP API timeout changed to: %s", tm)
		}
		if cmd.Flag("set-tcp-server-address").Changed {
			address, _ := cmd.Flags().GetString("set-tcp-server-address")
			_, _, err := net.SplitHostPort(address)
			if err != nil {
				jww.FEEDBACK.Println("Invalid TCP API server address given")
				os.Exit(1)
			}
			tcpApi.Address = address
			jww.FEEDBACK.Printf("TCP API server address changed to: %s", address)
		}
		if cmd.Flag("set-tcp-buffer-size").Changed {
			bufsz, err := cmd.Flags().GetInt("set-tcp-buffer-size")
			if err != nil || bufsz < 0 {
				jww.FEEDBACK.Println("Invalid buffer size given")
				return
			}
			tcpApi.BufferSize = bufsz
			jww.FEEDBACK.Println("TCP API buffer size changed to: %d", bufsz)
		}

		err = db.UpdateUserInfo(info)
		if err != nil {
			jww.ERROR.Println(err)
		}
		err = db.UpdateHttpApiConfig(httpApi)
		if err != nil {
			jww.ERROR.Println(err)
		}
		err = db.UpdateTcpApiConfig(tcpApi)
		if err != nil {
			jww.ERROR.Println(err)
		}
	},
}

func init() {
	configCmd.PersistentFlags().StringP("set-token", "t", "", "Set user token")
	configCmd.PersistentFlags().StringP("set-default-namespace", "n", "", "Default namespace")
	configCmd.PersistentFlags().StringP("set-http-server-address", "H", dbconfig.DefaultHTTPServer, "HTTP API server address")
	configCmd.PersistentFlags().Duration("set-http-server-timeout", dbconfig.DefaultHTTPTimeout, "HTTP API calls timeout")
	configCmd.PersistentFlags().StringP("set-tcp-server-address", "T", dbconfig.DefaultTCPServer, "TCP API server address")
	configCmd.PersistentFlags().Int("set-tcp-buffer-size", dbconfig.DefaultBufferSize, "TCP API buffer size")
	RootCmd.AddCommand(configCmd)
}
