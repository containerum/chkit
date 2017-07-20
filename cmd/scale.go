package cmd

import (
	"os"
	"strconv"

	"chkit-v2/chlib"
	"chkit-v2/helpers"
	"github.com/spf13/cobra"
)

var scaleCmdName string
var scaleCmdCount int

var scaleCmd = &cobra.Command{
	Use:   "scale KIND NAME COUNT",
	Short: "Change replicas count for object",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			np.FEEDBACK.Println("Invalid argument count")
			cmd.Usage()
			os.Exit(1)
		}
		switch args[0] {
		case "deployments", "deployment", "deploy":
			break
		default:
			np.FEEDBACK.Println("Invalid KIND. Choose from ('deployments', 'deployment', 'deploy')")
			cmd.Usage()
			os.Exit(1)
		}
		scaleCmdName = args[1]
		var err error
		scaleCmdCount, err = strconv.Atoi(args[2])
		if err != nil || scaleCmdCount <= 0 {
			np.FEEDBACK.Println("COUNT must be positive integer")
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
		ns, _ := getCmd.PersistentFlags().GetString("namespace")
		np.FEEDBACK.Println("scale...")
		err = client.Scale(scaleCmdName, scaleCmdCount, ns)
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
	scaleCmd.PersistentFlags().StringP("namespace", "n", "", "Namespace")
	RootCmd.AddCommand(scaleCmd)
}
