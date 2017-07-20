package cmd

import (
	"os"

	"chkit-v2/chlib"
	"chkit-v2/helpers"

	"github.com/spf13/cobra"
)

var deleteCmdKind, deleteCmdFile, deleteCmdName string

var deleteCmd = &cobra.Command{
	Use:   "delete (KIND NAME| --file -f FILE)",
	Short: "Remove object from namespace",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			np.FEEDBACK.Println("KIND or file not specified")
			cmd.Usage()
			os.Exit(1)
		}
		switch args[0] {
		case "--file", "-f":
			deleteCmdFile = args[1]
		case "po", "pods", "pod":
			deleteCmdKind = chlib.KindPods
		case "deployments", "deployment", "deploy":
			deleteCmdKind = chlib.KindDeployments
		case "service", "services", "svc":
			deleteCmdKind = chlib.KindService
		case "ns", "namespaces", "namespace":
			deleteCmdKind = chlib.KindNamespaces
		default:
			np.FEEDBACK.Println("Invalid KIND (choose from 'po', 'pods', 'pod', 'deployments', 'deployment', 'deploy', 'service', 'services', 'svc', 'ns', 'namespaces', 'namespace') or file")
			cmd.Usage()
			os.Exit(1)
		}
		if len(args) >= 2 && args[1][0] != '-' {
			deleteCmdName = args[1]
		} else {
			np.FEEDBACK.Println("NAME is not specified")
			cmd.Usage()
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := chlib.NewClient(db, helpers.CurrentClientVersion, helpers.UuidV4(), np)
		if err != nil {
			np.ERROR.Println(err)
			return
		}
		nameSpace, _ := cmd.Flags().GetString("namespace")
		np.FEEDBACK.Print("delete...")
		err = client.Delete(deleteCmdKind, deleteCmdName, nameSpace, false)
		if err != nil {
			np.FEEDBACK.Println("ERROR")
			np.ERROR.Println(err)
		} else {
			np.FEEDBACK.Println("OK")
		}
	},
}

func init() {
	deleteCmd.PersistentFlags().StringP("namespace", "n", "", "Namespace")
	RootCmd.AddCommand(deleteCmd)
}
