package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/kfeofantov/chkit-v2/chlib"
	"github.com/kfeofantov/chkit-v2/helpers"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var createCmd = &cobra.Command{
	Use:   "create (-f FILE | --file FILE)",
	Short: "Create object using JSON file",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flag("file").Changed || cmd.Flag("file").Value.String() == "" {
			jww.FEEDBACK.Println("File argument must be specified")
			cmd.Usage()
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.Flags().GetString("file")
		file, err := os.Open(filePath)
		if err != nil {
			jww.ERROR.Printf("file open: %s", err)
			return
		}
		content, err := ioutil.ReadAll(file)
		if err != nil {
			jww.ERROR.Printf("file read: %s", err)
			return
		}
		var jsonContent chlib.GenericJson
		err = json.Unmarshal(content, &jsonContent)
		if err != nil {
			jww.ERROR.Println("JSON parse: %s", err)
			return
		}
		client, err := chlib.NewClient(helpers.CurrentClientVersion, helpers.UuidV4())
		if err != nil {
			jww.ERROR.Println(err)
			return
		}
		jww.FEEDBACK.Print("create... ")
		_, err = client.Create(jsonContent)
		if err != nil {
			jww.FEEDBACK.Println("ERROR")
			jww.ERROR.Println(err)
		} else {
			jww.FEEDBACK.Println("OK")
		}
	},
}

func init() {
	createCmd.PersistentFlags().StringP("file", "f", "", "path to JSON file")
	RootCmd.AddCommand(createCmd)
}
