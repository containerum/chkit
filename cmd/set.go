package cmd

import (
	"os"
	"strings"

	"chkit-v2/chlib"
	"chkit-v2/helpers"

	"github.com/spf13/cobra"
)

var setCmdDeploy, setCmdContainer, setCmdParameter, setCmdValue string

var setCmd = &cobra.Command{
	Use:        "set KIND DEPLOY [CONTAINER] PARAMETER=VALUE",
	Short:      "Change one of parameters in Deployment",
	ValidArgs:  []string{chlib.KindDeployments},
	ArgAliases: []string{"deployments", "deployment", "deploy"},
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			cmd.Usage()
			os.Exit(1)
		}
		switch args[0] {
		case "deployments", "deployment", "deploy":
			break
		default:
			np.FEEDBACK.Printf("Invalid KIND. Choose from (%s)\n", strings.Join(cmd.ArgAliases, ", "))
			cmd.Usage()
			os.Exit(1)
		}
		setCmdDeploy = args[1]
		fieldValuePos := 2
		if len(args) == 4 {
			setCmdContainer = args[2]
			fieldValuePos++
		}
		if kv := strings.Split(args[fieldValuePos], "="); len(kv) == 2 {
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
		_, err = client.Set(setCmdDeploy, setCmdContainer, setCmdParameter, setCmdValue, ns)
		if err == nil {
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
