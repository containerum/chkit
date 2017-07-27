package cmd

import (
	"os"

	"chkit-v2/chlib"
	"chkit-v2/chlib/requestresults"
	"chkit-v2/helpers"
	"strings"

	"fmt"

	"regexp"

	"github.com/spf13/cobra"
)

var getCmdFile, getCmdKind, getCmdName string

var validOutputFormats = []string{"json", "yaml", "pretty"}

var getCmd = &cobra.Command{
	Use:        "get (KIND [NAME]| --file -f FILE)",
	Short:      "Show info about pod(s), service(s), namespace(s), deployment(s)",
	ValidArgs:  []string{chlib.KindPods, chlib.KindDeployments, chlib.KindNamespaces, chlib.KindService, "--file", "-f"},
	ArgAliases: []string{"po", "pods", "pod", "deployments", "deployment", "deploy", "service", "services", "svc", "ns", "namespaces", "namespace"},
	PreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("output").Changed {
			switch val, _ := cmd.Flags().GetString("output"); val {
			case "json", "yaml", "pretty", "list":
			default:
				np.FEEDBACK.Printf("Invalid output format. Choose from (%s)", strings.Join(validOutputFormats, ", "))
				cmd.Usage()
				os.Exit(1)
			}
		}
		if len(args) == 0 {
			np.FEEDBACK.Println("KIND or file not specified")
			cmd.Usage()
			os.Exit(1)
		}
		switch args[0] {
		case "po", "pods", "pod":
			getCmdKind = chlib.KindPods
		case "deployments", "deployment", "deploy":
			getCmdKind = chlib.KindDeployments
		case "service", "services", "svc":
			getCmdKind = chlib.KindService
		case "ns", "namespaces", "namespace":
			getCmdKind = chlib.KindNamespaces
		default:
			if cmd.Flag("file").Changed {
				getCmdFile, _ = cmd.Flags().GetString("file")
			} else {
				np.FEEDBACK.Printf("Invalid KIND. Choose from (%s)\n", strings.Join(cmd.ArgAliases, ", "))
				cmd.Usage()
				os.Exit(1)
			}
		}
		if len(args) >= 1 && getCmdFile == "" {
			if len(args) >= 2 {
				getCmdName = args[1]
			}
		} else {
			np.FEEDBACK.Println("KIND or FILE is not specified")
			cmd.Usage()
			os.Exit(1)
		}
		if getCmdName != "" && !regexp.MustCompile(chlib.ObjectNameRegex).MatchString(getCmdName) {
			np.FEEDBACK.Println("NAME is invalid")
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
		var jsonContent []chlib.GenericJson
		if getCmdFile != "" {
			err = chlib.LoadJsonFromFile(getCmdFile, &jsonContent)
		} else {
			nameSpace, _ := cmd.Flags().GetString("namespace")
			jsonContent, err = chlib.GetCmdRequestJson(client, getCmdKind, getCmdName, nameSpace)
		}
		if err != nil {
			np.ERROR.Printf("json receive error: %s\n", err)
			return
		}
		switch format, _ := cmd.Flags().GetString("output"); format {
		case "pretty":
			fieldToSort, _ := cmd.Flags().GetString("sort-by")
			p, err := requestresults.ProcessResponse(jsonContent, strings.ToUpper(fieldToSort), np)
			if err != nil {
				break
			}
			err = p.Print()
		case "json":
			err = chlib.JsonPrettyPrint(jsonContent, np)
		case "yaml":
			err = chlib.YamlPrint(jsonContent, np)
		}
		if err != nil {
			np.ERROR.Println(err)
		}
	},
}

func init() {
	getCmd.PersistentFlags().StringP("output", "o", "pretty", fmt.Sprintf("Output format: %s", strings.Join(validOutputFormats, ", ")))
	cobra.MarkFlagCustom(getCmd.PersistentFlags(), "output", "__chkit_get_outformat")
	getCmd.PersistentFlags().StringP("sort-by", "s", "NAME", "Sort by field. Used only if list printed")
	cobra.MarkFlagCustom(getCmd.PersistentFlags(), "sort-by", "__chkit_get_sort_columns")
	getCmd.PersistentFlags().StringP("namespace", "n", "", "Namespace")
	cobra.MarkFlagCustom(getCmd.PersistentFlags(), "namespace", "__chkit_namespaces_list")
	getCmd.PersistentFlags().StringP("file", "f", "", "JSON file generated on object creation")
	cobra.MarkFlagFilename(getCmd.PersistentFlags(), "file", "json")
	RootCmd.AddCommand(getCmd)
}
