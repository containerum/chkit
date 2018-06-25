package container

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/ninedraft/boxofstuff/str"
)

func componentVolumes(cont *container.Container, volumes str.Vector) *activekit.MenuItem {
	var containerVolumes = append([]model.ContainerVolume{}, cont.VolumeMounts...)
	return &activekit.MenuItem{
		Label: "Edit volume mounts",
		Action: func() error {
			for exit := false; !exit; {
				var volumeMounts = make(activekit.MenuItems, 0, len(containerVolumes))
				for _, vol := range containerVolumes {
					volumeMounts = append(volumeMounts, componentVolume(&vol, volumes, nil))
				}
				volumeMounts = volumeMounts.Append(&activekit.MenuItem{
					Label: "Add container",
					Action: func() error {
						var vol = model.ContainerVolume{}
						var ok = false
						componentVolume(&vol, volumes, &ok).Action()
						if ok {
							containerVolumes = append(containerVolumes, vol)
						}
						return nil
					},
				})
				(&activekit.Menu{
					Title: "Container -> Volume mounts\nType 'del VOLUME_NUMBER' to delete volume mount",
					Items: volumeMounts.
						Append(&activekit.MenuItem{
							Label: "Confirm",
							Action: func() error {
								exit = true
								cont.VolumeMounts = containerVolumes
								return nil
							},
						}, &activekit.MenuItem{
							Label: "Drop changes",
							Action: func() error {
								exit = true
								return nil
							},
						}),
					CustomOptionHandler: func(query string) error {
						var tokens = str.Fields(query).
							Map(strings.ToLower).
							Map(strings.TrimSpace)
						if tokens.Len() != 2 || (tokens.Len() > 0 && tokens[0] != "del") {
							fmt.Printf("invalid query %q: expecting 'del VOLUME_NUMBER'\n", query)
						}
						var ind, err = strconv.ParseUint(tokens[1], 10, 16)
						if err != nil {
							fmt.Printf("unable to parse volume number %q \n", tokens[1])
							return nil
						}
						if ind < 1 || int(ind) > len(containerVolumes) {
							fmt.Printf("Index of volume mount must be between %d and %d\n", ind, len(containerVolumes))
							return nil
						}
						var volLabel = str.Vector{containerVolumes[ind-1].Name, containerVolumes[ind-1].MountPath}.Join(" ")
						if activekit.YesNo("Are you sure you want to delete container %q?", volLabel) {
							containerVolumes = append(containerVolumes[:ind-1], containerVolumes[ind:]...)
						}
						return nil
					},
				}).Run()
			}
			return nil
		},
	}
}

func componentVolume(oldVol *model.ContainerVolume, volumes str.Vector, ok *bool) *activekit.MenuItem {
	return &activekit.MenuItem{
		Label: fmt.Sprintf("Edit volume mount %s",
			str.Vector{oldVol.Name, oldVol.MountPath}.Join(" ")),
		Action: func() error {
			var vol = *oldVol
			for exit := false; !exit; {
				var label = str.Vector{vol.Name, vol.MountPath, "empty volume mount"}.
					Filter(str.Longer(0)).
					Head(2).
					Join(" ")
				(&activekit.Menu{
					Title: "Container -> Volume mounts -> " + label,
					Items: activekit.MenuItems{
						{
							Label: fmt.Sprintf("Edit volume : %s", vol.Name),
							Action: func() error {
								(&activekit.Menu{
									Title: "Container -> Volume mounts -> " +
										label + " -> Select volume",
									Items: activekit.StringSelector(volumes, func(s string) error {
										vol.Name = s
										if vol.MountPath == "" {
											vol.MountPath = "/mnt/" + s
										}
										return nil
									}).Append(&activekit.MenuItem{
										Label: "Return to previous menu, leave " +
											str.Vector{vol.Name, "empty"}.FirstNonEmpty(),
									}),
								}).Run()
								return nil
							},
						},
						{
							Label: fmt.Sprintf("Edit path   : %s", vol.MountPath),
							Action: func() error {
								for {
									var volumePath = activekit.Promt("Type volume path, hit Enter to leave %s: ",
										str.Vector{vol.MountPath, "empty"}.FirstNonEmpty())
									volumePath = strings.TrimSpace(volumePath)
									if volumePath != "" {
										if !path.IsAbs(volumePath) {
											fmt.Printf("Mount path must be absolute!\n")
											continue
										} else {
											vol.MountPath = volumePath
										}
									}
									return nil
								}
							},
						},
					}.Append(&activekit.MenuItem{
						Label: "Confirm",
						Action: func() error {
							exit = true
							*oldVol = vol
							if ok != nil {
								*ok = true
							}
							return nil
						},
					}, &activekit.MenuItem{
						Label: "Return to previous menu, drop changes",
						Action: func() error {
							exit = true
							return nil
						},
					}),
				}).Run()
			}
			return nil
		},
	}
}
