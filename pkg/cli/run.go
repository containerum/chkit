package cli

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/containerum/chkit/pkg/cli/clisetup"
	"github.com/containerum/chkit/pkg/cli/login"
	"github.com/containerum/chkit/pkg/cli/mode"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/configdir"
	"github.com/containerum/chkit/pkg/configuration"
	. "github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// Version -- chkit version
	Version = "3.0.0-alpha"
)

const (
	// FlagConfigFile -- context config data key
	FlagConfigFile = "config"
	// FlagAPIaddr -- API address context key
	FlagAPIaddr = "apiaddr"
)

var (
	// ErrFatalError -- unrecoverable fatal error
	ErrFatalError chkitErrors.Err = "fatal error"
)

var runContext = struct {
	ConfigFile    string
	APIaddr       string
	Username      string
	Pass          string
	DebugRequests bool
}{}

var Root = &cobra.Command{
	Use:     "chkit",
	Short:   "chkit is a containerum.io terminal client",
	Version: Version,
	PersistentPreRun: func(*cobra.Command, []string) {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC1123,
		})
		logFile := path.Join(configdir.LogDir(), configuration.LogFileName())
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			logrus.Fatalf("error while creating log file: %v", err)
		}
		logrus.SetOutput(file)
	},
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Debugf("loading config")
		if err := configuration.LoadConfig(); err != nil {
			logrus.WithError(err).Errorf("unable to load config")
			fmt.Printf("Unable to load config :(\n")
			return
		}
		logrus.Debugf("running setup")
		err := clisetup.SetupConfig()
		switch {
		case err == nil:
			// pass
		case clisetup.ErrInvalidUserInfo.Match(err):
			logrus.Debugf("invalid user information")
			logrus.Debugf("running login")

			if err := login.Login(); err != nil {
				logrus.WithError(err).Errorf("unable to login")
				fmt.Printf("Unable to login: %v", err)
				return
			}
		default:
			logrus.WithError(ErrFatalError.Wrap(err)).Errorf("fatal erorr while login")
			angel.Angel(err)
			return
		}
		logrus.Debugf("client initialisation")
		if err := clisetup.SetupClient(); err != nil {
			logrus.WithError(err).Errorf("unable to init client")
			angel.Angel(err)
		}
		logrus.Debugf("saving tokens")
		if err := configuration.SaveTokens(Context.Tokens); err != nil {
			logrus.WithError(err).Errorf("unable to save tokens")
			fmt.Printf("Unable to save tokens!")
			return
		}

		logrus.Debugf("getting user namespaces list")
		list, err := Context.Client.GetNamespaceList()
		if err != nil {
			logrus.WithError(err).Errorf("unable to get user namespace list")
			fmt.Printf("Unable to get default namespace\n")
		}
		if len(list) == 0 {
			fmt.Printf("You have no namespaces!\n")
		}
		logrus.Infof("Hello, %q!", Context.Client.Username)
		if err := mainActivity(); err != nil {
			logrus.Fatalf("error in main activity: %v", err)
		}
	},
}

func init() {
	Context.Client.APIaddr = mode.API_ADDR
	Root.AddCommand(
		login.Command,
	)
	Root.PersistentFlags().StringVarP(&Context.Namespace, "namespace", "n", Context.Namespace, "")
	Root.PersistentFlags().BoolVarP(&Context.Quiet, "quiet", "q", Context.Quiet, "quiet mode")
}
