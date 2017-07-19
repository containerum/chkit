package cmd

import (
	"github.com/spf13/cobra"
	"chkit-v2/chlib/dbconfig"
	"os"
	"chkit-v2/chlib"
	jww "github.com/spf13/jwalterweatherman"
)

var db *dbconfig.ConfigDB

var debug bool

//RootCmd main cmd entrypoint
var RootCmd = &cobra.Command{
	Use: "chkit",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		db, err = dbconfig.OpenOrCreate(chlib.ConfigFile)
		if err != nil {
			jww.ERROR.Println(err)
			os.Exit(1)
		}
		debug, _ = cmd.Flags().GetBool("debug")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			cmd.Usage()
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		db.Close()
	},
}

func init() {
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "turn on debugging messages")
}