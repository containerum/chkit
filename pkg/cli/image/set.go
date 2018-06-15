package image

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/model/image"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/spf13/cobra"
)

var setAliases = []string{"imgs", "img", "im", "images"}

func Set(ctx *context.Context) *cobra.Command {
	force := false
	deplName := ""
	img := model.UpdateImage{}
	command := &cobra.Command{
		Use:     "image",
		Aliases: setAliases,
		Short:   "set container image in specific deployment",
		Long: `Sets container image in specific deployment.
If deployment contains only one container, then uses that container by default.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := prerun.PreRun(ctx); err != nil {
				activekit.Attention(err.Error())
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if force {
				if img.Container == "" {
					depl, err := ctx.Client.GetDeployment(ctx.Namespace.ID, deplName)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					img.Container = depl.Containers.Names()[0]
				}
				if err := ctx.Client.SetContainerImage(ctx.Namespace.ID, deplName, img); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return
			}
			config := image.Config{}

			if !cmd.Flag("deployment").Changed {
				deplList, err := ctx.Client.GetDeploymentList(ctx.Namespace.ID)
				if err != nil {
					activekit.Attention(err.Error())
				}
				var menu []*activekit.MenuItem
				for _, depl := range deplList {
					menu = append(menu, &activekit.MenuItem{
						Label: depl.Name,
						Action: func(depl deployment.Deployment) func() error {
							return func() error {
								deplName = depl.Name
								return nil
							}
						}(depl),
					})
				}
				(&activekit.Menu{
					Title: "Select deployment",
					Items: menu,
				}).Run()
			}

			if cmd.Flag("container").Changed {
				config.UpdateImage.Container = img.Container
			} else {
				depl, err := ctx.Client.GetDeployment(ctx.Namespace.ID, deplName)
				if err != nil {
					activekit.Attention(err.Error())
					os.Exit(1)
				}
				config.Containers = depl.Containers
			}
			if cmd.Flag("image").Changed {
				config.UpdateImage.Image = img.Image
			}
			if force {
				if err := ctx.Client.SetContainerImage(ctx.Namespace.ID, deplName, img); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return
			}
			img = image.Wizard(config)
			for exit := false; !exit; {
				(&activekit.Menu{
					Items: []*activekit.MenuItem{
						{
							Label: "Update image on server",
							Action: func() error {
								if !activekit.YesNo(fmt.Sprintf("Do you really want to update image of container %q?", img.Container)) {
									return nil
								}
								if err := ctx.Client.SetContainerImage(ctx.Namespace.ID, deplName, img); err != nil {
									activekit.Attention(err.Error())
									return nil
								}
								fmt.Printf("Image of container %q updated\n", img.Container)
								return nil
							},
						},
						{
							Label: "Edit image",
							Action: func() error {
								config.UpdateImage = img
								img = image.Wizard(config)
								return nil
							},
						},
						{
							Label: "Exit",
							Action: func() error {
								exit = true
								return nil
							},
						},
					},
				}).Run()
			}
		},
	}
	command.PersistentFlags().
		StringVarP(&deplName, "deployment", "d", "", "deployment label")
	command.PersistentFlags().
		BoolVarP(&force, "force", "f", false, "suppress confirmation")
	command.PersistentFlags().
		StringVarP(&img.Image, "image", "i", "", "new image")
	command.PersistentFlags().
		StringVarP(&img.Container, "container", "c", "", "container label")
	return command
}
