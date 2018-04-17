package cliserv

import (
	"strings"

	"fmt"

	"os"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var deleteServiceConfig = struct {
	Force bool
}{}

var Delete = &cobra.Command{
	Use:     "service",
	Aliases: aliases,
	Short:   "call to delete service in specific namespace",
	Long:    "deletes service in namespace. Aliases: " + strings.Join(aliases, ", "),
	Example: "chkit delete service service_label [-n namespace]",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Debugf("running command delete service")
		if len(args) == 0 {
			logrus.Debugf("showing help")
			cmd.Help()
			return
		}
		svcName := args[0]
		if !deleteServiceConfig.Force {
			if yes, _ := activekit.Yes(fmt.Sprintf("Do you really want delete service %q?", svcName)); !yes {
				return
			}
		}
		logrus.Debugf("deleting service %q from %q", svcName)
		err := func() error {
			return context.GlobalContext.Client.DeleteService(context.GlobalContext.Namespace, svcName)
		}()
		if err != nil {
			logrus.WithError(err).Debugf("error while deleting service")
			fmt.Printf("Unable to delete service %q :(\n%v", svcName, err)
			os.Exit(1)
		}
		fmt.Printf("OK\n")
		return
	},
}

func init() {
	Delete.PersistentFlags().
		BoolVarP(&deleteServiceConfig.Force, "force", "f", false, "force delete without confirmation")
}
