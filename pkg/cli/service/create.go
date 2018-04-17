package cliserv

import (
	"strings"

	"fmt"
	"os"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/service/servactive"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Create = &cobra.Command{
	Use:     "service",
	Aliases: aliases,
	Short:   "create service",
	Long:    "create service for provided pod in provided namespace. Aliases: " + strings.Join(aliases, ", "),
	Run: func(cmd *cobra.Command, args []string) {
		depList, err := context.GlobalContext.Client.GetDeploymentList(context.GlobalContext.Namespace)
		if err != nil {
			logrus.WithError(err).Errorf("unable to get deployment list")
			fmt.Println("Unable to get deployment list :(")
		}
		service, err := servactive.Wizard(servactive.ConstructorConfig{
			Deployments: depList.Names(),
		})
		if err != nil {
			logrus.WithError(err).Errorf("unable to create service")
			fmt.Println("Unable to create service :(")
			os.Exit(1)
		}
		fmt.Println(service.RenderTable())
	},
}
