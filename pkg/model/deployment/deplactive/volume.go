package deplactive

import (
	"fmt"
	"path/filepath"

	"io"

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
			if componentContainerVolume(vol, volumes).Action() == nil {
				cont.VolumeMounts = append(cont.VolumeMounts, *vol)
			}
			return nil
		},
	})
}

func componentContainerVolume(oldVolume *model.ContainerVolume, volumes []string) *activekit.MenuItem {
	var volume = *oldVolume
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
				_, err := (&activekit.Menu{
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
								_, err := (&activekit.Menu{
									Title: "Select volume",
									Items: activekit.StringSelector(volumes, func(volumeName string) error {
										volume.Name = volumeName
										if volume.MountPath == "" {
											volume.MountPath = "/mnt/" + volume.Name
										}
										return nil
									}).Append(&activekit.MenuItem{
										Label: fmt.Sprintf("Return to previous menu, leave %s",
											activekit.OrString(volume.Name, "empty")),
									}),
								}).Run()
								return err
							},
						},
						{
							Label: "Confirm",
							Action: func() error {
								if volume.Name != "" && volume.MountPath != "" {
									exit = true
									*oldVolume = volume
									return nil
								}
								if volume.Name == "" {
									fmt.Printf("Volume name must be non-empty!\n")
								}
								if volume.MountPath == "" {
									fmt.Printf("Volume mountpath must be non-empty!\n")
								}
								return nil
							},
						},
						{
							Label: "Drop changes, return to previous menu",
							Action: func() error {
								exit = true
								return io.EOF
							},
						},
					},
				}).Run()
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}
