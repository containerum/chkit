package cmd

import (
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"os"
	"chkit-v2/chlib"
	"chkit-v2/helpers"
)

var restartCmdName string

var restartCmd = &cobra.Command{
	Use:   "restart NAME",
	Short: "Restart pods by deploy name",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args)<1 {
			jww.FEEDBACK.Println("Deployment name must be specified")
			cmd.Usage()
			os.Exit(1)
		}
		restartCmdName = args[0]
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := chlib.NewClient(db, helpers.CurrentClientVersion, helpers.UuidV4())
		if err != nil {
			jww.ERROR.Println(err)
			return
		}
		nameSpace, _ := cmd.Flags().GetString("namespace")
		jww.FEEDBACK.Print("restart...")
		err = client.Delete(chlib.KindDeployments, restartCmdName, nameSpace, true)
		if err!=nil {
			jww.FEEDBACK.Println("ERROR")
			jww.ERROR.Println(err)
			os.Exit(1)
		} else {
			jww.FEEDBACK.Println("OK")
		}
	},
}

func init() {
	restartCmd.PersistentFlags().StringP("namespace", "n", "","Namespace")
	RootCmd.AddCommand(restartCmd)
}
