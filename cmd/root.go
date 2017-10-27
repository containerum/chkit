package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/containerum/chkit/chlib"
	"github.com/containerum/chkit/chlib/dbconfig"

	"fmt"
	"strings"

	"github.com/containerum/chkit/chlib/requestresults"

	"github.com/containerum/chkit/helpers"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var np *jww.Notepad

var client *chlib.Client

var bashCompletionFunc = fmt.Sprintf(`__chkit_get_outformat()
{
	COMPREPLY=( $(compgen -W "json yaml pretty" -- "${cur}") )
}

__chkit_get_object_list()
{
	local prog
	prog="${COMP_WORDS[0]}"
	ret=$(${prog} get $1 | cut -d '|' -f 2 | grep -vE '(^\+)|(NAME)')
	code=$?
	echo "${ret}"
	return ${code}
}

__chkit_containers_in_deploy()
{
	local prog
	prog="${COMP_WORDS[0]}"
	ret=$(${prog} get deploy $1 | sed -e '1,/Containers:/d' | sed '/Ports/,/ImagePullPolicy/d' | sed 's/\t//g')
	code=$?
	echo "${ret}"
	return ${code}
}

__chkit_namespaces_list()
{
	__chkit_get_object_list namespace
}

__chkit_get_sort_columns()
{
	local cur cmd cols
	cur="${COMP_WORDS[COMP_CWORD]}"
	cmd="${COMP_WORDS[2]}"
	case "${cmd}" in
		"deployments" | "deployment" | "deploy" )
			cols="%s"
		;;
		"po" | "pods" | "pod" )
			cols="%s"
		;;
		"service" | "services" | "svc" )
			cols="%s"
		;;
		"ns" | "namespaces" | "namespace" )
			cols="%s"
		;;
	esac
	COMPREPLY=( $(compgen -W "${cols}" -- "${cur}") )
}

__custom_func()
{
	local cur prev list
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[COMP_CWORD-1]}"
	case "${prev}" in
		"po" | "pods" | "pod" | "deployments" | "deployment" | "deploy" | "service" | "services" | "svc" | "ns" | "namespaces" | "namespace" )
			list="$(__chkit_get_object_list ${prev})"
			if [[ $? == 0 ]]; then
				COMPREPLY=( $(compgen -W "${list}" -- "${cur}") )
			fi
		;;
		* )
			if [[ ${last_command} == "chkit_set" && $COMP_CWORD == 4 ]]; then
				list="$(__chkit_containers_in_deploy ${COMP_WORDS[3]})"
				if [[ $? == 0 ]]; then
					COMPREPLY=( $(compgen -W "${list}" -- "${cur}") )
				fi
			fi
		;;
	esac
}`, strings.Join(requestresults.DeployColumns, " "),
	strings.Join(requestresults.PodColumns, " "),
	strings.Join(requestresults.ServiceColumns, " "),
	strings.Join(requestresults.NamespaceColumns, " "))

func exitOnErr(err error) {
	if err != nil {
		np.ERROR.Println(err)
		os.Exit(1)
	}
}

func saveUserSettings(cfg dbconfig.UserInfo) {
	np.FEEDBACK.Println("Saving settings")
	db, err := dbconfig.OpenOrCreate(chlib.ConfigFile, np)
	exitOnErr(err)
	exitOnErr(db.UpdateUserInfo(*client.UserConfig))
}

//RootCmd main cmd entrypoint
var RootCmd = &cobra.Command{
	Use: "chkit",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug, _ := cmd.Flags().GetBool("debug"); debug {
			np = jww.NewNotepad(jww.LevelDebug, jww.LevelDebug, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
		} else {
			np = jww.NewNotepad(jww.LevelInfo, jww.LevelInfo, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
		}
		db, err := dbconfig.OpenOrCreate(chlib.ConfigFile, np)
		exitOnErr(err)
		apiConfig, err := db.GetHttpApiConfig()
		exitOnErr(err)
		tcpApiConfig, err := db.GetTcpApiConfig()
		exitOnErr(err)
		userConfig, err := db.GetUserInfo()
		exitOnErr(err)
		uuid := helpers.UuidV4()

		client = &chlib.Client{
			ApiHandler:    &chlib.HttpApiHandler{Config: &apiConfig, UserInfo: &userConfig, Np: np, Channel: uuid},
			TcpApiHandler: &chlib.TcpApiHandler{Config: &tcpApiConfig, UserInfo: &userConfig, Np: np, Channel: uuid},
			UserConfig:    &userConfig,
		}
		exitOnErr(db.Close())
	},
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			cmd.Usage()
		}
	},
	BashCompletionFunction: bashCompletionFunc,
}

func init() {
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "turn on debugging messages")
}
