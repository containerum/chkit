package cmd

import (
	"os"
	"strconv"
	"strings"

	"github.com/containerum/chkit/chlib"

	"github.com/spf13/cobra"
)

var exposeCmdName string

var exposeCmdPorts []chlib.Port

var exposeCmd = &cobra.Command{
	Use:        "expose KIND NAME (-p --ports PORTS)",
	Short:      "Create Service and set output port list",
	ValidArgs:  []string{chlib.KindDeployments},
	ArgAliases: []string{"deployments", "deployment", "deploy"},
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
			np.FEEDBACK.Printf("Invalid KIND. Choose from (%s)\n", strings.Join(cmd.ArgAliases, ", "))
			cmd.Usage()
			os.Exit(1)
		}
		if !chlib.ObjectNameRegex.MatchString(args[1]) {
			np.FEEDBACK.Println("Invalid NAME specified")
			cmd.Usage()
			os.Exit(1)
		}
		exposeCmdName = args[1]
		portMatcher := chlib.PortRegex
		ports, _ := cmd.Flags().GetStringSlice("ports")
		for _, portStr := range ports {
			if !portMatcher.MatchString(portStr) {
				np.FEEDBACK.Println("Invalid PORT format. Must be PORTNAME:TARGETPORT[:PROTOCOL] or PORTNAME:TARGETPORT:PORT[:PROTOCOL]. Protocol is TCP or UDP")
				cmd.Usage()
				os.Exit(1)
			}
			subm := portMatcher.FindStringSubmatch(portStr)
			if !chlib.PortNameRegex.MatchString(subm[1]) {
				np.FEEDBACK.Println("Invalid port name")
				cmd.Usage()
				os.Exit(1)
			}
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
		nameSpace, _ := cmd.Flags().GetString("namespace")
		np.FEEDBACK.Print("expose...")
		_, err := client.Expose(exposeCmdName, exposeCmdPorts, nameSpace)
		if err != nil {
			np.FEEDBACK.Println("ERROR")
			np.ERROR.Println(err)
		} else {
			np.FEEDBACK.Println("OK")
		}
	},
}

func init() {
	exposeCmd.PersistentFlags().StringSliceP("ports", "p", []string{}, "Port list. Format PORTNAME:TARGETPORT[:PROTOCOL] or PORTNAME:TARGETPORT:PORT[:PROTOCOL]")
	cobra.MarkFlagRequired(exposeCmd.PersistentFlags(), "ports")
	exposeCmd.PersistentFlags().StringP("namespace", "n", "", "Namespace")
	cobra.MarkFlagCustom(exposeCmd.PersistentFlags(), "namespace", "__chkit_namespaces_list")
	RootCmd.AddCommand(exposeCmd)
}
