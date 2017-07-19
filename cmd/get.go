package cmd

import (
	"encoding/json"
	"os"

	"chkit-v2/chlib"
	"chkit-v2/chlib/requestresults"
	"chkit-v2/helpers"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"gopkg.in/yaml.v2"
)

var getCmdFile, getCmdKind, getCmdName string

var getCmd = &cobra.Command{
	Use:   "get (KIND [NAME]| --file -f FILE)",
	Short: "Show info about pod(s), service(s), namespace(s), deployment(s)",
	PreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("output").Changed {
			switch val, _ := cmd.Flags().GetString("output"); val {
			case "json", "yaml", "pretty":
			default:
				jww.FEEDBACK.Println("output must be json, yaml or pretty")
				cmd.Usage()
				os.Exit(1)
			}
		}
		if len(args) == 0 {
			jww.FEEDBACK.Println("KIND or file not specified")
			cmd.Usage()
			os.Exit(1)
		}
		switch args[0] {
		case "--file", "-f":
			getCmdFile = args[1]
			return
		case "po", "pods", "pod":
			getCmdKind = chlib.KindPods
		case "deployments", "deployment", "deploy":
			getCmdKind = chlib.KindDeployments
		case "service", "services", "svc":
			getCmdKind = chlib.KindService
		case "ns", "namespaces", "namespace":
			getCmdKind = chlib.KindNamespaces
		default:
			jww.FEEDBACK.Println("Invalid KIND (choose from 'po', 'pods', 'pod', 'deployments', 'deployment', 'deploy', 'service', 'services', 'svc', 'ns', 'namespaces', 'namespace') or file")
			cmd.Usage()
			os.Exit(1)
		}
		if len(args) >= 2 && args[1][0] != '-' {
			getCmdName = args[1]
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := chlib.NewClient(db, helpers.CurrentClientVersion, helpers.UuidV4())
		if err != nil {
			jww.ERROR.Println(err)
			return
		}
		var jsonContent []chlib.GenericJson
		if getCmdFile != "" {
			err = chlib.LoadJsonFromFile(getCmdFile, &jsonContent)
		} else {
			nameSpace, _ := cmd.Flags().GetString("namespace")
			jsonContent, err = chlib.GetCmdRequestJson(client, getCmdKind, getCmdName, nameSpace)
		}
		if err != nil {
			jww.ERROR.Printf("json receive error: %s\n", err)
			return
		}
		switch format, _ := cmd.Flags().GetString("output"); format {
		case "pretty":
			fieldToSort, _ := cmd.Flags().GetString("sort-by")
			p, err := requestresults.ProcessResponse(jsonContent, fieldToSort)
			if err != nil {
				break
			}
			err = p.Print()
		case "json":
			var b []byte
			b, err = json.MarshalIndent(jsonContent, "", "    ")
			jww.FEEDBACK.Printf("%s\n", b)
		case "yaml":
			var b []byte
			b, err = yaml.Marshal(jsonContent)
			jww.FEEDBACK.Printf("%s\n", b)
		}
		if err != nil {
			jww.ERROR.Println(err)
		}
	},
}

func init() {
	getCmd.PersistentFlags().StringP("output", "o", "pretty", "Output format: json, yaml, pretty")
	getCmd.PersistentFlags().StringP("sort-by", "s", "NAME", "Sort by field. Used only if list printed")
	getCmd.PersistentFlags().StringP("namespace", "n", "", "Namespace")
	getCmd.PersistentFlags().StringP("file", "f", "", "JSON file generated on object creation")
	RootCmd.AddCommand(getCmd)
}
