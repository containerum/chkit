package clinamespace

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Delete(ctx *context.Context) *cobra.Command {
	var deleteNamespaceConfig = struct {
		Force bool
	}{}
	command := &cobra.Command{
		Use:     "namespace",
		Short:   "call to delete namespace",
		Long:    "delete namespace deletes namespace with name, provided by first arg. Aliases: " + strings.Join(aliases, ", "),
		Example: "chkit delete namespace $ID",
		Aliases: aliases,
		Run: func(command *cobra.Command, args []string) {
			logrus.WithFields(logrus.Fields{
				"command": "delete namespace",
			}).Debugf("getting namespace data")
			if len(args) == 0 {
				logrus.Debugf("showing help")
				command.Help()
				return
			}
			ns := args[0]
			err := func() error {
				if !deleteNamespaceConfig.Force {
					yes, _ := activekit.Yes(fmt.Sprintf("Do you want to delete namespace %q?", ns))
					if !yes {
						return nil
					}
				}
				logrus.Debugf("deleting namespace %q", ns)
				return ctx.Client.DeleteNamespace(ns)
			}()
			if err != nil {
				logrus.WithError(err).Errorf("unable to delete namespace %q", ns)
				fmt.Printf("Unable to delete namespace :(\n%v", err)
				return
			}
			fmt.Printf("Namespace %q deleted\n", ns)
		},
	}
	command.PersistentFlags().BoolVarP(&deleteNamespaceConfig.Force, "force", "f", false, "force delete without confirmation")
	return command
}
