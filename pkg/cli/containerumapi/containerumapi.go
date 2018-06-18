package containerumapi

import (
	"fmt"
	"net/url"
	"os"

	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Set(ctx *context.Context) *cobra.Command {
	var flags struct {
		AllowSelfSignedCerts bool `desc:""`
	}
	command := &cobra.Command{
		Use:     "containerum-api",
		Short:   "Set Containerum API URL",
		Aliases: aliases,
		PreRun: prerun.PreRunFunc(ctx, prerun.Config{
			DoNotRunLoginOnIncompatibleConfig: true,
			SetupClient:                       false,
			AllowInvalidConfig:                true,
		}),
		Run: func(cmd *cobra.Command, args []string) {
			var logger = ctx.Log.Command("set containerum-api")
			logger.Debugf("START set containerum-api")
			defer logger.Debugf("END set containerum-api")
			logger.StructFields(flags)
			if len(args) != 1 && !flags.AllowSelfSignedCerts {
				logger.Debugf("invalid flags and args combination, showing help")
				cmd.Help()
				os.Exit(1)
			} else if len(args) == 1 {
				logger.Debugf("validating API URL %q", args[0])
				api, err := url.Parse(args[0])
				if err != nil {
					logger.WithError(err).Errorf("invalid API URL")
					fmt.Println(err)
					os.Exit(1)
				}
				ctx.Client.APIaddr = api.String()
			}
			ctx.AllowSelfSignedTLS = flags.AllowSelfSignedCerts
			ctx.Changed = true
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
