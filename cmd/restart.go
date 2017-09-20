package cmd

import (
	"os"

	"github.com/containerum/chkit.v2/chlib"
	"github.com/containerum/chkit.v2/helpers"
	"github.com/spf13/cobra"
)

var restartCmdName string

var restartCmd = &cobra.Command{
	Use:   "restart NAME",
	Short: "Restart pods by deploy name",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			np.FEEDBACK.Println("Deployment name must be specified")
			cmd.Usage()
			os.Exit(1)
		}
		if !chlib.ObjectNameRegex.MatchString(args[0]) {
			np.FEEDBACK.Println("Invalid NAME specified")
			cmd.Usage()
			os.Exit(1)
		}
		restartCmdName = args[0]
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := chlib.NewClient(db, helpers.CurrentClientVersion, helpers.UuidV4(), np)
		if err != nil {
			np.ERROR.Println(err)
			return
		}
		nameSpace, _ := cmd.Flags().GetString("namespace")
		np.FEEDBACK.Print("restart...")
		_, err = client.Delete(chlib.KindDeployments, restartCmdName, nameSpace, true)
		if err != nil {
			np.FEEDBACK.Println("ERROR")
			np.ERROR.Println(err)
			os.Exit(1)
		} else {
			np.FEEDBACK.Println("OK")
		}
	},
}

func init() {
	restartCmd.PersistentFlags().StringP("namespace", "n", "", "Namespace")
	cobra.MarkFlagCustom(restartCmd.PersistentFlags(), "namespace", "__chkit_namespaces_list")
	RootCmd.AddCommand(restartCmd)
}
