package cliserv

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/service"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Delete(ctx *context.Context) *cobra.Command {
	var deleteServiceConfig = struct {
		Force bool
	}{}
	command := &cobra.Command{
		Use:     "service",
		Aliases: aliases,
		Short:   "delete service in specific namespace",
		Long:    "Delete service in namespace.",
		Example: "chkit delete service service_label [-n namespace]",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Debugf("running command delete service")
			switch len(args) {
			case 0:
				list, err := ctx.Client.GetServiceList(ctx.Namespace.ID)
				if err != nil {
					logrus.WithError(err).Errorf("unable to get service list")
					activekit.Attention(err.Error())
					return
				}
				var menu []*activekit.MenuItem
				for _, srv := range list {
					menu = append(menu, &activekit.MenuItem{
						Label: srv.Name,
						Action: func(srv service.Service) func() error {
							return func() error {
								if yes, _ := activekit.Yes(fmt.Sprintf("Do you really want delete service %q?", srv.Name)); !yes {
									return nil
								}
								logrus.Debugf("deleting service %q from %q", srv.Name)
								err := ctx.Client.DeleteService(ctx.Namespace.ID, srv.Name)
								if err != nil {
									logrus.WithError(err).Debugf("error while deleting service")
									fmt.Printf("Unable to delete service %q :(\n%v", srv.Name, err)
									os.Exit(1)
								}
								fmt.Printf("OK\n")
								return nil
							}
						}(srv),
					})
				}
				(&activekit.Menu{
					Title: "Which service do you want to delete?",
					Items: append(menu, []*activekit.MenuItem{
						{
							Label: "Exit",
						},
					}...),
				}).Run()
			case 1:
				svcName := args[0]
				if !deleteServiceConfig.Force {
					if yes, _ := activekit.Yes(fmt.Sprintf("Do you really want delete service %q?", svcName)); !yes {
						return
					}
				}
				logrus.Debugf("deleting service %q from %q", svcName)
				err := ctx.Client.DeleteService(ctx.Namespace.ID, svcName)
				if err != nil {
					logrus.WithError(err).Debugf("error while deleting service")
					fmt.Printf("Unable to delete service %q :(\n%v", svcName, err)
					os.Exit(1)
				}
				fmt.Printf("OK\n")
				return
			default:
				logrus.Debugf("showing help")
				cmd.Help()
				return
			}
		},
	}
	command.PersistentFlags().
		BoolVarP(&deleteServiceConfig.Force, "force", "f", false, "force delete without confirmation")
	return command
}
