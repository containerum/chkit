package containerumapi

import (
	"fmt"

	"github.com/containerum/chkit/pkg/context"
	"github.com/spf13/cobra"
)

var aliases = []string{"api", "current-api", "api-addr", "API"}

func Get(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:     "containerum-api",
		Aliases: aliases,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(ctx.Client.APIaddr)
		},
	}
	return command
}
