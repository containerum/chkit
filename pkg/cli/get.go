package cli

import (
	"github.com/containerum/chkit/pkg/cli/deployment"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/configuration"
	. "github.com/containerum/chkit/pkg/context"
	"github.com/spf13/cobra"
)

var Get = &cobra.Command{
	Use: "get",
	PersistentPreRun: func(command *cobra.Command, args []string) {
		prerun.PreRun()
	},
	Run: func(command *cobra.Command, args []string) {
		command.Help()
	},
	PersistentPostRun: func(command *cobra.Command, args []string) {
		if Context.Changed {
			configuration.SaveConfig()
		}
	},
}

func init() {
	Get.AddCommand(
		clideployment.Get,
	)
}
