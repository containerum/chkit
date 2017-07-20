package cmd

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"chkit-v2/chlib"
	"chkit-v2/helpers"

	"github.com/spf13/cobra"
)

const portRegex = `^(\D+):(\d+)(:(\d+))?(:(TCP|UDP))?$`

var exposeCmdName string

var exposeCmdPorts []chlib.Port

var exposeCmd = &cobra.Command{
	Use:   "expose KIND NAME (-p --ports PORTS)",
	Short: "Create Service and set output port list",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			np.FEEDBACK.Println("Invalid argument count")
			cmd.Usage()
			os.Exit(1)
		}

		switch args[0] {
		case "deployments", "deployment", "deploy":
			break
		default:
			np.FEEDBACK.Println("Invalid KIND. Choose from ('deployments', 'deployment', 'deploy')")
			cmd.Usage()
			os.Exit(1)
		}
		exposeCmdName = args[1]
		portMatcher := regexp.MustCompile(portRegex)
		portsStr, _ := cmd.Flags().GetString("ports")
		for _, portStr := range strings.Split(portsStr, ",") {
			if !portMatcher.MatchString(portStr) {
				np.FEEDBACK.Println("Invalid PORT format. Must be PORTNAME:TARGETPORT[:PROTOCOL] or PORTNAME:TARGETPORT:PORT[:PROTOCOL]. Protocol is TCP or UDP")
				cmd.Usage()
				os.Exit(1)
			}
			subm := portMatcher.FindStringSubmatch(portStr)
			portName := subm[1]
			targetPort, err := strconv.Atoi(subm[2])
			if err != nil || targetPort < 0 || targetPort > 65535 {
				np.FEEDBACK.Println("Invalid target port")
				cmd.Usage()
				os.Exit(1)
			}
			var port int
			if subm[4] != "" {
				port, err = strconv.Atoi(subm[4])
				if err != nil || port < 0 || port > 65535 {
					np.FEEDBACK.Println("Invalid port")
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
		client, err := chlib.NewClient(db, helpers.CurrentClientVersion, helpers.UuidV4(), np)
		if err != nil {
			np.ERROR.Println(err)
			return
		}
		nameSpace, _ := cmd.Flags().GetString("namespace")
		np.FEEDBACK.Print("expose...")
		_, err = client.Expose(exposeCmdName, exposeCmdPorts, nameSpace)
		if err != nil {
			np.FEEDBACK.Println("ERROR")
			np.ERROR.Println(err)
		} else {
			np.FEEDBACK.Println("OK")
		}
	},
}

func init() {
	exposeCmd.PersistentFlags().StringP("ports", "p", "", "Port list. Format PORTNAME:TARGETPORT[:PROTOCOL] or PORTNAME:TARGETPORT:PORT[:PROTOCOL], split with \",\"")
	exposeCmd.PersistentFlags().StringP("namespace", "n", "", "Namespace")
	RootCmd.AddCommand(exposeCmd)
}
