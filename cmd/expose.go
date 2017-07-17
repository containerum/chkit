package cmd

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/kfeofantov/chkit-v2/chlib"
	"github.com/kfeofantov/chkit-v2/helpers"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

const portRegex = `^(\D+):(\d+)(:(\d+))?(:(TCP|UDP))?$`

var exposeCmdName string

var exposeCmdPorts []chlib.Port

var exposeCmd = &cobra.Command{
	Use:   "expose KIND NAME (-p --ports PORTS)",
	Short: "Create Service and set output port list",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			jww.FEEDBACK.Println("Invalid argument count")
			cmd.Usage()
			os.Exit(1)
		}

		switch args[0] {
		case "deployments", "deployment", "deploy":
			break
		default:
			jww.FEEDBACK.Println("Invalid KIND. Choose from ('deployments', 'deployment', 'deploy')")
			cmd.Usage()
			os.Exit(1)
		}
		exposeCmdName = args[1]
		portMatcher := regexp.MustCompile(portRegex)
		portsStr, _ := cmd.Flags().GetString("ports")
		for _, portStr := range strings.Split(portsStr, ",") {
			if !portMatcher.MatchString(portStr) {
				jww.FEEDBACK.Println("Invalid PORT format. Must be PORTNAME:TARGETPORT[:PROTOCOL] or PORTNAME:TARGETPORT:PORT[:PROTOCOL]. Protocol is TCP or UDP")
				cmd.Usage()
				os.Exit(1)
			}
			subm := portMatcher.FindStringSubmatch(portStr)
			portName := subm[1]
			targetPort, err := strconv.Atoi(subm[2])
			if err != nil || targetPort < 0 || targetPort > 65535 {
				jww.FEEDBACK.Println("Invalid target port")
				cmd.Usage()
				os.Exit(1)
			}
			var port int
			if subm[4] != "" {
				port, err = strconv.Atoi(subm[4])
				if err != nil || port < 0 || port > 65535 {
					jww.FEEDBACK.Println("Invalid port")
					cmd.Usage()
					os.Exit(1)
				}
			}
			proto := subm[6]
			if proto == "" {
				proto = chlib.DefaultProto
			}
			exposeCmdPorts = append(exposeCmdPorts, chlib.Port{
				TargetPort: targetPort,
				Port:       port,
				Name:       portName,
				Protocol:   proto,
			})
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := chlib.NewClient(helpers.CurrentClientVersion, helpers.UuidV4())
		if err != nil {
			jww.ERROR.Println(err)
			return
		}
		nameSpace, _ := cmd.Flags().GetString("namespace")
		jww.FEEDBACK.Print("expose...")
		_, err = client.Expose(exposeCmdName, exposeCmdPorts, nameSpace)
		if err != nil {
			jww.FEEDBACK.Println("ERROR")
			jww.ERROR.Println(err)
		} else {
			jww.FEEDBACK.Println("OK")
		}
	},
}

func init() {
	cfg, err := chlib.GetUserInfo()
	if err != nil {
		panic(err)
	}
	exposeCmd.PersistentFlags().StringP("ports", "p", "", "Port list. Format PORTNAME:TARGETPORT[:PROTOCOL] or PORTNAME:TARGETPORT:PORT[:PROTOCOL], split with \",\"")
	exposeCmd.PersistentFlags().StringP("namespace", "n", cfg.Namespace, "Namespace")
	RootCmd.AddCommand(exposeCmd)
}
