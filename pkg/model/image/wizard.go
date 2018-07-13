package image

import (
	"fmt"

	"strings"

	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/validation"
	kubeModel "github.com/containerum/kube-client/pkg/model"
)

type Config struct {
	Containers  []container.Container
	UpdateImage kubeModel.UpdateImage
}

func Wizard(config Config) kubeModel.UpdateImage {
	updImage := config.UpdateImage
	oldImage := updImage
	if updImage.Container == "" &&
		len(config.Containers) == 1 {
		updImage.Container = config.Containers[0].Name
	}
	var ok bool
	for exit := false; !exit; {
		(&activekit.Menu{
			Title: "Update image data",
			Items: []*activekit.MenuItem{
				{
					Label: fmt.Sprintf("Set image     : %s",
						activekit.OrString(updImage.Image, "none (required)")),
					Action: func() error {
						img := strings.TrimSpace(activekit.Promt("Type image: "))
						if img == "" {
							return nil
						}
						if err := validation.ValidateImageName(img); err != nil {
							activekit.Attention(fmt.Sprintf("Invalid image %q", img))
						}
						updImage.Image = img
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Set container : %s",
						activekit.OrString(updImage.Container, "none (required")),
					Action: func() error {
						var menu []*activekit.MenuItem
						for _, cont := range config.Containers {
							menu = append(menu, &activekit.MenuItem{
								Label: fmt.Sprintf("%s [%s]", cont.Name, cont.Image),
								Action: func(cont string) func() error {
									return func() error {
										updImage.Container = cont
										return nil
									}
								}(cont.Name),
							})
						}
						(&activekit.Menu{
							Title: "Which container do you want update?",
							Items: append(menu, []*activekit.MenuItem{
								{
									Label: "Set custom container",
									Action: func() error {
										contLabel := strings.TrimSpace(activekit.Promt("Type container label: "))
										if contLabel == "" {
											return nil
										}
										if err := validation.ValidateContainerName(contLabel); err != nil {
											activekit.Attention(fmt.Sprintf("Invalid container name"))
										}
										updImage.Container = contLabel
										return nil
									},
								},
							}...),
						}).Run()
						return nil
					},
				},
				{
					Label: "Confirm",
					Action: func() error {
						if err := ValidateImage(updImage); err != nil {
							activekit.Attention(err.Error())
							return nil
						}
						ok = true
						exit = true
						return nil
					},
				},
				{
					Label: "<-",
					Action: func() error {
						exit = true
						ok = false
						return nil
					},
				},
			},
		}).Run()
	}
	if ok {
		return updImage
	}
	return oldImage
}
