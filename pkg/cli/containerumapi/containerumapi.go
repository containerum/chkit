package containerumapi

import (
	"net/url"

	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/ferr"
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
			InitClient:             prerun.DoNotInitClient,
			RunLoginOnMissingCreds: false,
		}),
		Run: func(cmd *cobra.Command, args []string) {
			var logger = ctx.Log.Command("set containerum-api")
			logger.Debugf("START set containerum-api")
			defer logger.Debugf("END set containerum-api")
			logger.StructFields(flags)
			if len(args) != 1 && !flags.AllowSelfSignedCerts {
				logger.Debugf("invalid flags and args combination, showing help")
				cmd.Help()
				ctx.Exit(1)
			} else if len(args) == 1 {
				logger.Debugf("validating API URL %q", args[0])
				api, err := url.Parse(args[0])
				if err != nil {
					logger.WithError(err).Errorf("invalid API URL")
					ferr.Println(err)
					ctx.Exit(1)
				}
				ctx.GetClient().APIaddr = api.String()
			}
			ctx.SetSelfSignedTLSRule(flags.AllowSelfSignedCerts)
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
