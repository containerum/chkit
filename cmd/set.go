package cmd

import (
	"os"
	"strings"

	"chkit-v2/chlib"
	"chkit-v2/helpers"
	"github.com/spf13/cobra"
)

var setCmdField, setCmdParameter, setCmdValue, setCmdContainer string

var setCmd = &cobra.Command{
	Use:   "set FIELD TYPE CONTAINER PARAMETER=VALUE",
	Short: "Change one of parameters in Deployment",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) < 4 {
			cmd.Usage()
			os.Exit(1)
		}
		switch args[0] {
		case "image":
		default:
			np.FEEDBACK.Println("Invalid field name")
			os.Exit(1)
		}
		setCmdField = args[0]
		switch args[1] {
		case "deployments", "deployment", "deploy":
			break
		default:
			np.FEEDBACK.Println("Invalid KIND. Choose from ('deployments', 'deployment', 'deploy')")
			os.Exit(1)
		}
		setCmdContainer = args[2]
		if kv := strings.Split(args[3], "="); len(kv) == 2 {
			setCmdParameter = kv[0]
			setCmdValue = kv[1]
		} else {
			np.FEEDBACK.Println("Invalid parameter syntax")
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
		np.FEEDBACK.Print("set... ")
		_, err = client.Set(setCmdField, setCmdContainer, setCmdValue, ns)
		if err != nil {
			np.FEEDBACK.Println("OK")
		} else {
			np.FEEDBACK.Println("ERROR")
			np.ERROR.Println(err)
		}
	},
}

func init() {
	setCmd.PersistentFlags().StringP("namespace", "n", "", "Namespace")
	RootCmd.AddCommand(setCmd)
}
