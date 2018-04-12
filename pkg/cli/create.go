package cli

import (
	"github.com/containerum/chkit/pkg/cli/deployment"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/spf13/cobra"
)

var Create = &cobra.Command{
	Use: "create",
	PersistentPreRun: func(command *cobra.Command, args []string) {
		prerun.PreRun()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	Create.AddCommand(
		clideployment.Create,
	)
}
