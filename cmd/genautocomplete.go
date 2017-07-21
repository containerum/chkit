package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var genautocompleteCmd = &cobra.Command{
	Use:   "genautocomplete",
	Short: "Generate bash autocompletion parameters",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		fileName, _ := cmd.Flags().GetString("file")
		if fileName != "-" {
			err = RootCmd.GenBashCompletionFile(fileName)
		} else {
			err = RootCmd.GenBashCompletion(os.Stdout)
		}
		if err != nil {
			np.ERROR.Printf("Bash autocompletion generate: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	genautocompleteCmd.PersistentFlags().StringP("file", "f", "-", "File to write generated parameters, \"-\" for stdout")
	cobra.MarkFlagFilename(genautocompleteCmd.PersistentFlags(), "file")
	RootCmd.AddCommand(genautocompleteCmd)
}
