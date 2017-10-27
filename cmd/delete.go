package cmd

import (
	"os"

	"github.com/containerum/chkit/chlib"

	"strings"

	"github.com/spf13/cobra"
)

var deleteCmdKind, deleteCmdFile string
var deleteCmdNames []string

var deleteCmd = &cobra.Command{
	Use:        "delete (KIND NAME [NAME2]...| --file -f FILE)",
	Short:      "Remove object from namespace",
	ValidArgs:  []string{chlib.KindPods, chlib.KindDeployments, chlib.KindNamespaces, chlib.KindService, "--file", "-f"},
	ArgAliases: []string{"po", "pods", "pod", "deployments", "deployment", "deploy", "service", "services", "svc", "ns", "namespaces", "namespace"},
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			deleteCmdFile, _ = cmd.Flags().GetString("file")
			var obj chlib.CommonObject
			exitOnErr(chlib.LoadJsonFromFile(deleteCmdFile, &obj))
			deleteCmdKind = strings.ToLower(obj.Kind) + "s"
			deleteCmdNames = []string{obj.Metadata.Name}
		} else {
			if cmd.Flag("file").Changed {
				np.FEEDBACK.Println("Can`t use both file and KIND arguments")
				cmd.Usage()
				os.Exit(1)
			}
			switch args[0] {
			case "po", "pods", "pod":
				deleteCmdKind = chlib.KindPods
			case "deployments", "deployment", "deploy":
				deleteCmdKind = chlib.KindDeployments
			case "service", "services", "svc":
				deleteCmdKind = chlib.KindService
			case "ns", "namespaces", "namespace":
				deleteCmdKind = chlib.KindNamespaces
			default:
				np.FEEDBACK.Printf("Invalid KIND. Choose from (%s)\n", strings.Join(cmd.ArgAliases, ", "))
				cmd.Usage()
				os.Exit(1)
			}
			if len(args) >= 2 {
				for _, v := range args[1:] {
					if !chlib.ObjectNameRegex.MatchString(v) {
						break
					}
					deleteCmdNames = append(deleteCmdNames, v)
				}
			} else {
				np.FEEDBACK.Println("NAME is not specified or invalid")
				cmd.Usage()
				os.Exit(1)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		nameSpace, _ := cmd.Flags().GetString("namespace")
		for _, deleteCmdName := range deleteCmdNames {
			np.FEEDBACK.Printf("delete %s...", deleteCmdName)
			_, err := client.Delete(deleteCmdKind, deleteCmdName, nameSpace, false)
			if err != nil {
				np.FEEDBACK.Println("ERROR")
				np.ERROR.Println(err)
				return
			} else {
				np.FEEDBACK.Println("OK")
			}
		}
	},
}

func init() {
	deleteCmd.PersistentFlags().StringP("file", "f", "", "File generated at object creation")
	cobra.MarkFlagFilename(deleteCmd.PersistentFlags(), "file", "json")
	deleteCmd.PersistentFlags().StringP("namespace", "n", "", "Namespace")
	cobra.MarkFlagCustom(deleteCmd.PersistentFlags(), "namespace", "__chkit_namespaces_list")
	RootCmd.AddCommand(deleteCmd)
}
