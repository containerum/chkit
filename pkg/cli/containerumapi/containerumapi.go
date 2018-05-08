package containerumapi

import (
	"fmt"
	"net/url"
	"os"

	"github.com/containerum/chkit/pkg/context"
	"github.com/spf13/cobra"
)

func Set(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use: "containerum-api",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 && !cmd.Flags().Changed("allow-self-signed-certs") {
				cmd.Help()
				os.Exit(1)
			} else if len(args) == 1 {
				api, err := url.Parse(args[0])
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				ctx.Client.APIaddr = api.String()
			}
			ctx.AllowSelfSignedTLS, _ = cmd.Flags().GetBool("allow-self-signed-certs")
			ctx.Changed = true
		},
	}
	command.PersistentFlags().
		Bool("allow-self-signed-certs", false, "")
	return command
}
