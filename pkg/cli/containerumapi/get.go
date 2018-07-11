package containerumapi

import (
	"fmt"

	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/spf13/cobra"
)

var aliases = []string{"api", "current-api", "api-addr", "API"}

func Get(ctx *context.Context) *cobra.Command {
	configuration.LoadConfig(ctx)

	command := &cobra.Command{
		Use:     "containerum-api",
		Short:   "print Containerum API URL",
		Aliases: aliases,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(ctx.GetClient().APIaddr)
		},
	}
	return command
}
