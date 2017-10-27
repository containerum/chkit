package cmd

import (
	"encoding/base64"
	"net"
	"net/url"

	"os"

	"github.com/containerum/chkit/chlib"
	"github.com/containerum/chkit/chlib/dbconfig"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure chkit default values",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := dbconfig.OpenOrCreate(chlib.ConfigFile, np)
		exitOnErr(err)
		info, err := db.GetUserInfo()
		exitOnErr(err)
		httpApi, err := db.GetHttpApiConfig()
		exitOnErr(err)
		tcpApi, err := db.GetTcpApiConfig()
		exitOnErr(err)

		if cmd.Flags().NFlag() == 0 {
			np.FEEDBACK.Println("Token: ", info.Token)
			np.FEEDBACK.Println("Namespace: ", info.Namespace)
			np.FEEDBACK.Println("HTTP API")
			np.FEEDBACK.Println("\tServer: ", httpApi.Server)
			np.FEEDBACK.Println("\tTimeout: ", httpApi.Timeout)
			np.FEEDBACK.Println("TCP API")
			np.FEEDBACK.Printf("\tServer: %s", tcpApi.Address)
			np.FEEDBACK.Println("\tBuffer size: ", tcpApi.BufferSize)
			np.FEEDBACK.Println("\tTimeout: ", tcpApi.Timeout)
			return
		}

		if cmd.Flag("set-default-namespace").Changed {
			newNamespace := cmd.Flag("set-default-namespace").Value.String()
			if _, err := client.Get(chlib.KindNamespaces, newNamespace, ""); err != nil {
				np.ERROR.Println(err)
				os.Exit(1)
			}
			info.Namespace = newNamespace
			np.FEEDBACK.Printf("Namespace changed to: %s\n", info.Namespace)
		}
		if cmd.Flag("set-token").Changed {
			enteredToken := cmd.Flag("set-token").Value.String()
			if _, err := base64.StdEncoding.DecodeString(enteredToken); err != nil {
				np.FEEDBACK.Println("Invalid token given")
				os.Exit(1)
			}
			info.Token = enteredToken
			np.FEEDBACK.Printf("Token changed to: %s\n", info.Token)
		}
		if cmd.Flag("set-http-server-address").Changed {
			address := cmd.Flag("set-http-server-address").Value.String()
			if _, err := url.ParseRequestURI(address); err != nil {
				np.FEEDBACK.Printf("Invalid HTTP API server address given")
				os.Exit(1)
			}
			httpApi.Server = address
			np.FEEDBACK.Printf("HTTP API server address changed to: %s", address)
		}
		if cmd.Flag("set-http-server-timeout").Changed {
			tm, err := cmd.Flags().GetDuration("set-http-server-timeout")
			if err != nil {
				np.FEEDBACK.Printf("Invalid HTTP API timeout given")
				os.Exit(1)
			}
			httpApi.Timeout = tm
			np.FEEDBACK.Printf("HTTP API timeout changed to: %s", tm)
		}
		if cmd.Flag("set-tcp-server-timeout").Changed {
			tm, err := cmd.Flags().GetDuration("set-tcp-server-timeout")
			if err != nil {
				np.FEEDBACK.Printf("Invalid TCP API timeout given")
				os.Exit(1)
			}
			tcpApi.Timeout = tm
			np.FEEDBACK.Printf("TCP API timeout changed to: %s", tm)
		}
		if cmd.Flag("set-tcp-server-address").Changed {
			address, _ := cmd.Flags().GetString("set-tcp-server-address")
			_, _, err := net.SplitHostPort(address)
			if err != nil {
				np.FEEDBACK.Println("Invalid TCP API server address given")
				os.Exit(1)
			}
			tcpApi.Address = address
			np.FEEDBACK.Printf("TCP API server address changed to: %s", address)
		}
		if cmd.Flag("set-tcp-buffer-size").Changed {
			bufsz, err := cmd.Flags().GetInt("set-tcp-buffer-size")
			if err != nil || bufsz < 0 {
				np.FEEDBACK.Println("Invalid buffer size given")
				return
			}
			tcpApi.BufferSize = bufsz
			np.FEEDBACK.Println("TCP API buffer size changed to: %d", bufsz)
		}

		exitOnErr(db.UpdateUserInfo(info))
		exitOnErr(db.UpdateHttpApiConfig(httpApi))
		exitOnErr(db.UpdateTcpApiConfig(tcpApi))
		exitOnErr(db.Close())
	},
}

func init() {
	configCmd.PersistentFlags().StringP("set-token", "t", "", "Set user token")
	configCmd.PersistentFlags().StringP("set-default-namespace", "n", "", "Default namespace")
	cobra.MarkFlagCustom(configCmd.PersistentFlags(), "set-default-namespace", "__chkit_namespaces_list")
	configCmd.PersistentFlags().StringP("set-http-server-address", "H", dbconfig.DefaultHTTPServer, "HTTP API server address")
	configCmd.PersistentFlags().Duration("set-http-server-timeout", dbconfig.DefaultHTTPTimeout, "HTTP API calls timeout")
	configCmd.PersistentFlags().StringP("set-tcp-server-address", "T", dbconfig.DefaultTCPServer, "TCP API server address")
	configCmd.PersistentFlags().Int("set-tcp-buffer-size", dbconfig.DefaultBufferSize, "TCP API buffer size")
	configCmd.PersistentFlags().Duration("set-tcp-server-timeout", dbconfig.DefaultTCPTimeout, "TCP API calls timeout")
	RootCmd.AddCommand(configCmd)
}
