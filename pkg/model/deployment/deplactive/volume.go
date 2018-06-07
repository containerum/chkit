package deplactive

import (
	"fmt"
	"path/filepath"

	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/kube-client/pkg/model"
)

func componentEditContainerVolumes(cont *container.Container, volumes []string) activekit.MenuItems {
	var menuItems = make(activekit.MenuItems, 0, len(cont.VolumeMounts))
	for _, vol := range cont.VolumeMounts {
		menuItems = menuItems.Append(componentContainerVolume(&vol, volumes))
	}
	return menuItems.Append(&activekit.MenuItem{
		Label: "Mount volume",
		Action: func() error {
			var vol = &model.ContainerVolume{}
			componentContainerVolume(vol, volumes).Action()
			cont.VolumeMounts = append(cont.VolumeMounts, *vol)
			return nil
		},
	})
}

func componentContainerVolume(volume *model.ContainerVolume, volumes []string) *activekit.MenuItem {
	return &activekit.MenuItem{
		Label: fmt.Sprintf("Edit volume mount %s", func() string {
			if volume.MountPath != "" && volume.Name != "" {
				return volume.Name + " -> " + volume.MountPath
			}
			if volume.Name != "" {
				return volume.Name
			}
			if volume.MountPath != "" {
				return volume.MountPath
			}
			return ""
		}()),
		Action: func() error {
			for exit := false; !exit; {
				(&activekit.Menu{
					Title: "Deployment -> Container -> Volume",
					Items: activekit.MenuItems{
						{
							Label: fmt.Sprintf("Edit path %s",
								activekit.OrString(volume.MountPath, "undefined, required")),
							Action: activekit.HandleString(
								fmt.Sprintf("Type mount path (hit Enter to leave %s): ",
									activekit.OrString(volume.MountPath, "empty")),
								func(pathString string) bool {
									if !filepath.IsAbs(pathString) {
										fmt.Printf("Mount path must be absolute\n")
										return true
									}
									if pathString == "" {
										return true
									}
									volume.MountPath = pathString
									return false
								}),
						},
						{
							Label: fmt.Sprintf("Set volume : %s",
								activekit.OrString(volume.Name, "undefined. required")),
							Action: func() error {
								(&activekit.Menu{
									Title: "Select volume",
									Items: activekit.StringSelector(volumes, func(volumeName string) error {
										volume.Name = volumeName
										if volume.MountPath == "" {
											volume.MountPath = "/mnt/" + volume.Name
										}
										return nil
									}),
								}).Run()
								return nil
							},
						},
						{
							Label: "Confirm",
							Action: func() error {
								exit = true
								return nil
							},
						},
					},
				}).Run()
			}
			return nil
		},
	}
}
