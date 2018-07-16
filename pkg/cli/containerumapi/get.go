package containerumapi

import (
	"fmt"

	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/spf13/cobra"
)

var aliases = []string{"api", "current-api", "api-addr", "API"}

func Get(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:     "containerum-api",
		Short:   "print Containerum API URL",
		Aliases: aliases,
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := configuration.LoadConfig(ctx); err != nil {
				angel.Angel(ctx, err)
				ctx.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(ctx.GetClient().APIaddr)
		},
	}
	return command
}
