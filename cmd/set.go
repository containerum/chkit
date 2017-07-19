package cmd

import (
	"os"
	"strings"

	"chkit-v2/chlib"
	"chkit-v2/helpers"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
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
			jww.FEEDBACK.Println("Invalid field name")
			os.Exit(1)
		}
		setCmdField = args[0]
		switch args[1] {
		case "deployments", "deployment", "deploy":
			break
		default:
			jww.FEEDBACK.Println("Invalid KIND. Choose from ('deployments', 'deployment', 'deploy')")
			os.Exit(1)
		}
		setCmdContainer = args[2]
		if kv := strings.Split(args[3], "="); len(kv) == 2 {
			setCmdParameter = kv[0]
			setCmdValue = kv[1]
		} else {
			jww.FEEDBACK.Println("Invalid parameter syntax")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := chlib.NewClient(db, helpers.CurrentClientVersion, helpers.UuidV4())
		if err != nil {
			jww.ERROR.Println(err)
			return
		}
		ns, _ := getCmd.PersistentFlags().GetString("namespace")
		jww.FEEDBACK.Print("set... ")
		_, err = client.Set(setCmdField, setCmdContainer, setCmdValue, ns)
		if err != nil {
			jww.FEEDBACK.Println("OK")
		} else {
			jww.FEEDBACK.Println("ERROR")
			jww.ERROR.Println(err)
		}
	},
}

func init() {
	setCmd.PersistentFlags().StringP("namespace", "n", "", "Namespace")
	RootCmd.AddCommand(setCmd)
}
