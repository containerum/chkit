package cmd

import (
	"chkit-v2/chlib"
	"chkit-v2/chlib/dbconfig"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var db *dbconfig.ConfigDB

var np *jww.Notepad

//RootCmd main cmd entrypoint
var RootCmd = &cobra.Command{
	Use: "chkit",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug, _ := cmd.Flags().GetBool("debug"); debug {
			np = jww.NewNotepad(jww.LevelDebug, jww.LevelDebug, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
		} else {
			np = jww.NewNotepad(jww.LevelInfo, jww.LevelInfo, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
		}
		var err error
		db, err = dbconfig.OpenOrCreate(chlib.ConfigFile, np)
		if err != nil {
			np.ERROR.Println(err)
			os.Exit(1)
		}
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
