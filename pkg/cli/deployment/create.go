package clideployment

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/model/deployment/deplactive"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Create(ctx *context.Context) *cobra.Command {
	var flags deplactive.Flags
	command := &cobra.Command{
		Use:     "deployment",
		Aliases: aliases,
		Short:   "create new deployment",
		Long: "Creates new deployment.\n" +
			"Has an one-line mode, suitable for integration with other tools,\n" +
			"and an interactive wizard mod",
		Run: func(cmd *cobra.Command, args []string) {
			var depl deployment.Deployment
			if flags.File != "" {
				var depl, err = deplactive.FromFile(flags.File)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				flags = deplactive.FlagsFromDeployment(depl)
			} else {
				var err error
				depl, err = flags.Deployment()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
			deplactive.Fill(&depl)
			if flags.Force {
				if err := deplactive.ValidateDeployment(depl); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(depl.RenderYAML())
				if err := ctx.Client.CreateDeployment(ctx.Namespace.ID, depl); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println("OK")
				return
			}
			fmt.Println(depl.RenderTable())
		},
	}

	if err := gpflag.ParseTo(&flags, command.Flags()); err != nil {
		panic(err)
	}
	return command
}

func deploymentMenu(ctx *context.Context, depl deployment.Deployment, firstItem *activekit.MenuItem) *activekit.Menu {
	return &activekit.Menu{
		Items: []*activekit.MenuItem{
			firstItem,
			{
				Label: "Edit deployment",
				Action: func() error {
					var err error
					depl, err = deplactive.ReplaceWizard(deplactive.Config{
						Deployment: &depl,
					})
					if err != nil {
						logrus.WithError(err).Errorf("unable to create deployment")
						fmt.Println(err)
						os.Exit(1)
					}
					return nil
				},
			},
			{
				Label: "Print to terminal",
				Action: activekit.ActionWithErr(func() error {
					if data, err := depl.RenderYAML(); err != nil {
						return err
					} else {
						upBorders := strings.Repeat("_", text.Width(data))
						downBorders := strings.Repeat("_", text.Width(data))
						fmt.Printf("%s\n\n%s\n%s\n", upBorders, data, downBorders)
					}
					return nil
				}),
			},
			{
				Label: "Save to file",
				Action: func() error {
					filename, _ := activekit.AskLine("Print filename: ")
					data, err := depl.RenderJSON()
					if err != nil {
						return err
					}
					if err := ioutil.WriteFile(filename, []byte(data), os.ModePerm); err != nil {
						logrus.WithError(err).Errorf("unable to save deployment %q to file", depl.Name)
						fmt.Printf("Unable to save deployment to file :(\n%v", err)
						return nil
					}
					fmt.Printf("OK\n")
					return nil
				},
			},
			{
				Label: "Exit",
				Action: func() error {
					if yes, _ := activekit.Yes("Are you sure you want to exit?"); yes {
						os.Exit(0)
					}
					return nil
				},
			},
		},
	}
}
