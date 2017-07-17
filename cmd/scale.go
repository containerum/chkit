package cmd

import (
	"os"
	"strconv"

	"github.com/kfeofantov/chkit-v2/chlib"
	"github.com/kfeofantov/chkit-v2/helpers"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var scaleCmdName string
var scaleCmdCount int

var scaleCmd = &cobra.Command{
	Use:   "scale KIND NAME COUNT",
	Short: "Change replicas count for object",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			jww.FEEDBACK.Println("Invalid argument count")
			cmd.Usage()
			os.Exit(1)
		}
		switch args[0] {
		case "deployments", "deployment", "deploy":
			break
		default:
			jww.FEEDBACK.Println("Invalid KIND. Choose from ('deployments', 'deployment', 'deploy')")
			cmd.Usage()
			os.Exit(1)
		}
		scaleCmdName = args[1]
		var err error
		scaleCmdCount, err = strconv.Atoi(args[2])
		if err != nil || scaleCmdCount <= 0 {
			jww.FEEDBACK.Println("COUNT must be positive integer")
			cmd.Usage()
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := chlib.NewClient(helpers.CurrentClientVersion, helpers.UuidV4())
		if err != nil {
			jww.ERROR.Println(err)
			return
		}
		ns, _ := getCmd.PersistentFlags().GetString("namespace")
		jww.FEEDBACK.Println("scale...")
		err = client.Scale(scaleCmdName, scaleCmdCount, ns)
		if err != nil {
			jww.FEEDBACK.Println("ERROR")
			jww.ERROR.Println(err)
			os.Exit(1)
		} else {
			jww.FEEDBACK.Println("OK")
		}
	},
}

func init() {
	cfg, err := chlib.GetUserInfo()
	if err != nil {
		panic(err)
	}
	scaleCmd.PersistentFlags().StringP("namespace", "n", cfg.Namespace, "Namespace")
	RootCmd.AddCommand(scaleCmd)
}
