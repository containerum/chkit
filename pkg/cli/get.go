package cli

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/containerum/chkit/pkg/cli/configmap"
	"github.com/containerum/chkit/pkg/cli/containerumapi"
	"github.com/containerum/chkit/pkg/cli/deployment"
	"github.com/containerum/chkit/pkg/cli/ingress"
	"github.com/containerum/chkit/pkg/cli/namespace"
	"github.com/containerum/chkit/pkg/cli/pod"
	"github.com/containerum/chkit/pkg/cli/service"
	"github.com/containerum/chkit/pkg/cli/solution"
	"github.com/containerum/chkit/pkg/cli/user"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/spf13/cobra"
)

func Get(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:   "get",
		Short: "Get resource data",
		Run: func(command *cobra.Command, args []string) {
			command.Help()
		},
		PersistentPostRun: func(command *cobra.Command, args []string) {
			if ctx.Changed {
				if err := configuration.SyncConfig(ctx); err != nil {
					logrus.WithError(err).Errorf("unable to save config")
					fmt.Printf("Unable to save config: %v\n", err)
					return
				}
			}
			if err := configuration.SaveTokens(ctx, ctx.Client.Tokens); err != nil {
				logrus.WithError(err).Errorf("unable to save tokens")
				fmt.Printf("Unable to save tokens: %v\n", err)
				return
			}
		},
	}
	command.AddCommand(
		clideployment.Get(ctx),
		clinamespace.Get(ctx),
		clinamespace.GetAccess(ctx),
		cliserv.Get(ctx),
		clipod.Get(ctx),
		clingress.Get(ctx),
		cliuser.Get(ctx),
		clisolution.Get(ctx),
		containerumapi.Get(ctx),
		cliconfigmap.Get(ctx),
		&cobra.Command{
			Use:     "default-namespace",
			Short:   "print default",
			Aliases: []string{"default-ns", "def-ns"},
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Printf("%s\n", ctx.Namespace)
			},
		},
	)
	command.PersistentFlags().
		StringP("namespace", "n", ctx.Namespace.ID, "")
	return command
}
