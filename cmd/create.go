package cmd

import (
	"os"

	"github.com/containerum/chkit/chlib"
	"github.com/containerum/chkit/helpers"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create (-f FILE | --file FILE)",
	Short: "Create object using JSON file",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flag("file").Changed || cmd.Flag("file").Value.String() == "" {
			np.FEEDBACK.Println("File argument must be specified")
			cmd.Usage()
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.Flags().GetString("file")
		var jsonContent chlib.GenericJson
		err := chlib.LoadJsonFromFile(filePath, &jsonContent)
		if err != nil {
			np.ERROR.Println(err)
			return
		}
		client, err := chlib.NewClient(db, helpers.CurrentClientVersion, helpers.UuidV4(), np)
		if err != nil {
			np.ERROR.Println(err)
			return
		}
		np.FEEDBACK.Print("create... ")
		_, err = client.Create(jsonContent)
		if err != nil {
			np.FEEDBACK.Println("ERROR")
			np.ERROR.Println(err)
		} else {
			np.FEEDBACK.Println("OK")
		}
	},
}

func init() {
	createCmd.PersistentFlags().StringP("file", "f", "", "path to JSON file")
	cobra.MarkFlagRequired(createCmd.PersistentFlags(), "file")
	cobra.MarkFlagFilename(createCmd.PersistentFlags(), "file", "json")
	RootCmd.AddCommand(createCmd)
}
