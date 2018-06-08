package clideployment

import (
	"fmt"
	"os"

	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/model/deployment/deplactive"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Create(ctx *context.Context) *cobra.Command {
	var flags deplactive.Flags
	command := &cobra.Command{
		Use:     "deployment",
		Aliases: aliases,
		Short:   "create deployment",
		Long: "Creates new deployment.\n" +
			"Has an one-line mode, suitable for integration with other tools,\n" +
			"and an interactive wizard mod",
		Run: func(cmd *cobra.Command, args []string) {
			var depl deployment.Deployment
			var err error
			if flags.File != "" {
				depl, err = deplactive.FromFile(flags.File)
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
				if err := ctx.Client.CreateDeployment(ctx.Namespace.ID, depl); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Printf("Deployment %s created\n", depl.Name)
				return
			}
			configmapList, err := ctx.Client.GetConfigmapList(ctx.Namespace.ID)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(depl.RenderTable())
			depl = deplactive.Wizard{
				EditName:   true,
				Deployment: &depl,
				Configmaps: configmapList.Names(),
			}.Run()
			fmt.Println(depl.RenderTable())
			if !activekit.YesNo("Are you sure you want create deployment %s?", depl.Name) {
				(&activekit.Menu{
					Items: activekit.MenuItems{
						{
							Label: fmt.Sprintf("Save deployment %s to file", depl.Name),
							Action: func() error {
								for {
									var fname = activekit.Promt("Type filename (if ext is yaml or yml then file encodes as YAML, JSON by default): ")
									fname = strings.TrimSpace(fname)
									var err error
									var data string
									switch strings.ToLower(filepath.Ext(fname)) {
									case ".yaml", ".yml":
										fmt.Println("Encoding deployment as YAML")
										data, err = depl.RenderYAML()
									default:
										fmt.Println("Encoding deployment as JSON")
										data, err = depl.RenderJSON()
									}
									if err != nil {
										fmt.Println(err)
									}
									if err := ioutil.WriteFile(fname, []byte(data), os.ModePerm); err != nil {
										fmt.Println(err)
									}
								}
								return nil
							},
						},
					},
				}).Run()
				return
			}
			if err := ctx.Client.CreateDeployment(ctx.Namespace.ID, depl); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("Deployment %s created\n", depl.Name)
		},
	}
	if err := gpflag.ParseTo(&flags, command.Flags()); err != nil {
		angel.Angel(ctx, fmt.Errorf("it seems that the structure of the flags is set incorrectly: %v", err))
	}
	return command
}
