package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/containerum/chkit/pkg/client"
	"github.com/containerum/chkit/pkg/model"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

const (
	FlagAPIaddr    = "apiaddr"
	FlagConfigFile = "config"
)

var Configuration = struct {
	ConfigPath   string
	ConfigFile   string
	TokenFile    string
	ClientConfig model.Config
}{}

var ChkitClient client.ChkitClient
var notepad *jww.Notepad
var App = &cobra.Command{
	Use: "chkit",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug, _ := cmd.Flags().GetBool("debug"); debug {
			notepad = jww.NewNotepad(jww.LevelDebug,
				jww.LevelDebug,
				os.Stdout,
				ioutil.Discard,
				"", log.Ldate|log.Ltime)
		} else {
			notepad = jww.NewNotepad(jww.LevelInfo,
				jww.LevelInfo,
				os.Stdout,
				ioutil.Discard,
				"", log.Ldate|log.Ltime)
		}

	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	App.PersistentFlags().StringVar(&Configuration.ConfigFile,
		FlagConfigFile, "", "config file (default is "+Configuration.ConfigPath+"/containerum.yaml)")
	viper.BindPFlag(FlagConfigFile, App.PersistentFlags().Lookup(FlagConfigFile))

	App.PersistentFlags().StringVar(&Configuration.ClientConfig.APIaddr,
		FlagAPIaddr, "", "API address")

}
func initConfig() {
	viper.AddConfigPath(Configuration.ConfigPath)
	viper.SetConfigName("containerum")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
