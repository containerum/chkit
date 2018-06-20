package image

import (
	"fmt"
	"os"

	"strings"

	"errors"

	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/validation"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

var setAliases = []string{"imgs", "img", "im", "images"}

func Set(ctx *context.Context) *cobra.Command {
	var flags struct {
		Force      bool   `desc:"suppress confirmation" flag:"force f"`
		Deployment string `desc:"deployment name"`
		Container  string `desc:"container name"`
		Image      string `desc:"new image"`
	}
	var buildImage = func() (deployment string, image model.UpdateImage, err error) {
		var errs []string
		if validation.ValidateImageName(flags.Image) != nil {
			errs = append(errs, fmt.Sprintf("invalid image name %q", flags.Image))
		}
		if validation.ValidateContainerName(flags.Container) != nil {
			errs = append(errs, fmt.Sprintf("invalid container name %q", flags.Container))
		}
		if validation.ValidateLabel(flags.Deployment) != nil {
			errs = append(errs, fmt.Sprintf("invalid deployment name %q", flags.Deployment))
		}
		if len(errs) > 0 {
			return "", model.UpdateImage{}, errors.New(strings.Join(errs, "\n"))
		}
		return flags.Deployment, model.UpdateImage{
			Image:     flags.Image,
			Container: flags.Container,
		}, nil
	}
	command := &cobra.Command{
		Use:     "image",
		Aliases: setAliases,
		Short:   "Set container image for specific deployment.",
		Long: "Set container image for specific deployment\n" +
			"If a deployment contains only one container, the command will use that container by default.",
		PreRun: prerun.PreRunFunc(ctx),
		Run: func(cmd *cobra.Command, args []string) {
			var logger = ctx.Log.Command("set image")
			logger.Debugf("START")
			defer logger.Debugf("END")
			logger.StructFields(flags)
			if flags.Force {
				logger.Debugf("run command with force")

				if flags.Container == "" {
					var depl, err = ctx.Client.GetDeployment(ctx.Namespace.ID, flags.Deployment)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					if len(depl.Containers) == 1 {
						flags.Container = depl.Containers[0].Name
					}
				}

				var depl, image, err = buildImage()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if err := ctx.Client.SetContainerImage(ctx.Namespace.ID, depl, image); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println("OK")
				return
			}
			if flags.Deployment == "" {
				var deplList, err = ctx.Client.GetDeploymentList(ctx.Namespace.ID)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				(&activekit.Menu{
					Title: "Select deployment",
					Items: activekit.StringSelector(deplList.Names(), func(s string) error {
						flags.Deployment = s
						return nil
					}),
				}).Run()
			}
			if flags.Container == "" {
				var depl, err = ctx.Client.GetDeployment(ctx.Namespace.ID, flags.Deployment)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if len(depl.Containers) == 1 {
					flags.Container = depl.Containers[0].Name
				} else {
					(&activekit.Menu{
						Title: "Select container",
						Items: activekit.StringSelector(depl.Containers.Names(), func(s string) error {
							flags.Container = s
							return nil
						}),
					}).Run()
				}
			}
			if flags.Image == "" {
				for {
					var image = activekit.Promt("Type new image: ")
					image = strings.TrimSpace(image)
					if validation.ValidateImageName(image) != nil {
						fmt.Printf("%q is invalid image name!\n", image)
						continue
					}
					flags.Image = image
					break
				}
			}
			if activekit.YesNo("Are you sure you want to update image to %q of container %q in deployment %s/%s?",
				flags.Image, flags.Container, ctx.Namespace, flags.Deployment) {
				if err := ctx.Client.SetContainerImage(ctx.Namespace.ID, flags.Deployment, model.UpdateImage{
					Image:     flags.Image,
					Container: flags.Container,
				}); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
