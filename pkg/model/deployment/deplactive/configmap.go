package deplactive

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/kube-client/pkg/model"
)

func configmapsMenu(oldCm []model.ContainerVolume, configmaps []string) []model.ContainerVolume {
	var newCms = append(make([]model.ContainerVolume, 0, len(oldCm)), oldCm...)
	var items = make(activekit.MenuItems, 0, len(oldCm))
	for _, configmap := range newCms {
		items = items.Append(&activekit.MenuItem{
			Label: fmt.Sprintf("Edit %s %q", configmap.Name,
				activekit.OrString(configmap.MountPath, "/")),
			Action: func(config model.ContainerVolume) func() error {
				return func() error {

					return nil
				}
			}(configmap),
		})
	}
	for exit := false; !exit; {
		(&activekit.Menu{
			Title: `What's next? type "del OPTION" to delete configmap`,
			Items: items.Append(&activekit.MenuItem{
				Label: "Create configmap",
				Action: func() error {
					var configmapName string
					(&activekit.Menu{
						Title: "Select configmap",
						Items: activekit.SelectString(configmaps, func(s string) error {
							configmapName = s
							return nil
						}),
					}).Run()
					return nil
				},
			}),
			CustomOptionHandler: func(query string) error {
				var tokens = strings.Fields(query)
				switch {
				case len(tokens) == 2 && tokens[0] == "del":
					if i, err := strconv.Atoi(tokens[1]); err != nil && i > 0 && i <= len(oldCm) {
						i--
						newCms = append(newCms[:i], newCms[i+1:]...)
					} else if err != nil {
						fmt.Println(err)
					}
				}
				return nil
			},
		}).Run()
	}
	return newCms
}

func volumeNames(volumes []model.ContainerVolume) []string {
	var names = make([]string, 0, len(volumes))
	for _, volume := range volumes {
		names = append(names, volume.Name)
	}
	return names
}
