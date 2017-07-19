package cmd

import (
	"os"

	"chkit-v2/chlib"
	"chkit-v2/helpers"
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
		var jsonContent chlib.GenericJson
		err := chlib.LoadJsonFromFile(filePath, &jsonContent)
		if err != nil {
			jww.ERROR.Println(err)
			return
		}
		client, err := chlib.NewClient(db, helpers.CurrentClientVersion, helpers.UuidV4())
		if err != nil {
			jww.ERROR.Println(err)
			return
		}
		jww.FEEDBACK.Print("create... ")
		err = client.Create(jsonContent)
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
