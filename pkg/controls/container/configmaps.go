package container

import (
	"fmt"
	"path"
	"strings"

	"strconv"

	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/ninedraft/boxofstuff/str"
)

func componentConfigmaps(cont *container.Container, configmaps str.Vector) *activekit.MenuItem {
	var containerConfigs = append([]model.ContainerVolume{}, cont.ConfigMaps...)
	return &activekit.MenuItem{
		Label: "Edit configmaps",
		Action: func() error {
			for exit := false; !exit; {
				var configmapsMenuItems = make(activekit.MenuItems, 0, len(containerConfigs))
				for _, config := range containerConfigs {
					configmapsMenuItems = append(configmapsMenuItems, componentConfigmap(&config, configmaps, nil))
				}
				configmapsMenuItems = configmapsMenuItems.Append(&activekit.MenuItem{
					Label: "Mount configmap",
					Action: func() error {
						var vol = model.ContainerVolume{}
						var ok = false
						componentConfigmap(&vol, configmaps, &ok).Action()
						if ok {
							containerConfigs = append(containerConfigs, vol)
						}
						return nil
					},
				})
				(&activekit.Menu{
					Title: "Container -> Configmaps\nType 'del CONFIG_NUMBER' to delete volume mount",
					Items: configmapsMenuItems.
						Append(&activekit.MenuItem{
							Label: "Confirm",
							Action: func() error {
								exit = true
								cont.ConfigMaps = containerConfigs
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
							fmt.Printf("invalid query %q: expecting 'del CONFIG_NUMBER'\n", query)
						}
						var ind, err = strconv.ParseUint(tokens[1], 10, 16)
						if err != nil {
							fmt.Printf("unable to parse volume number %q \n", tokens[1])
							return nil
						}
						if ind < 1 || int(ind) > len(containerConfigs) {
							fmt.Printf("Index of volume mount must be between %d and %d\n", ind, len(containerConfigs))
							return nil
						}
						var volLabel = str.Vector{containerConfigs[ind-1].Name, containerConfigs[ind-1].MountPath}.Join(" ")
						if activekit.YesNo("Are you sure you want to delete container %q?", volLabel) {
							containerConfigs = append(containerConfigs[:ind-1], containerConfigs[ind:]...)
						}
						return nil
					},
				}).Run()
			}
			return nil
		},
	}
}

func componentConfigmap(oldConfig *model.ContainerVolume, configmaps str.Vector, ok *bool) *activekit.MenuItem {
	return &activekit.MenuItem{
		Label: fmt.Sprintf("Edit configmap %s",
			str.Vector{oldConfig.Name, oldConfig.MountPath}.Join(" ")),
		Action: func() error {
			var config = *oldConfig
			for exit := false; !exit; {
				var label = str.Vector{config.Name, config.MountPath, "empty configmap mount"}.
					Filter(str.Longer(0)).
					Head(2).
					Join(" ")
				(&activekit.Menu{
					Title: "Container -> Configmaps -> " + label,
					Items: activekit.MenuItems{
						{
							Label: fmt.Sprintf("Edit configmap : %s", config.Name),
							Action: func() error {
								(&activekit.Menu{
									Title: "Container -> Configmaps -> " +
										label + " -> Configmap",
									Items: activekit.StringSelector(configmaps, func(s string) error {
										config.Name = s
										if config.MountPath == "" {
											config.MountPath = "/etc/" + s
										}
										return nil
									}).Append(&activekit.MenuItem{
										Label: "Return to previous menu, leave " +
											str.Vector{config.Name, "empty"}.FirstNonEmpty(),
									}),
								}).Run()
								return nil
							},
						},
						{
							Label: fmt.Sprintf("Edit path      : %s", config.MountPath),
							Action: func() error {
								for {
									var volumePath = activekit.Promt("Type configmap path, hit Enter to leave %s: ",
										str.Vector{config.MountPath, "empty"}.FirstNonEmpty())
									volumePath = strings.TrimSpace(volumePath)
									if volumePath != "" {
										if !path.IsAbs(volumePath) {
											fmt.Printf("Mount path must be absolute!\n")
											continue
										} else {
											config.MountPath = volumePath
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
							*oldConfig = config
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
