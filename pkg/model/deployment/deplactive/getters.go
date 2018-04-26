package deplactive

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/containerum/chkit/pkg/util/validation"
)

func getContainers(conts []container.Container) []container.Container {
	containers := make([]container.Container, len(conts))
	copy(containers, conts)
	ok := true
	for exit := false; !exit; {
		containerMenuItems := make([]*activekit.MenuItem, 0, len(containers))
		for i, cont := range containers {
			containerMenuItems = append(containerMenuItems,
				&activekit.MenuItem{
					Label: fmt.Sprintf("Edit container %q", cont.Name),
					Action: func(i int, cont container.Container) func() error {
						return func() error {
							logrus.Debugf("editing container %q", containerMenuItems[i])
							edited, ok := getContainer(cont)
							if ok {
								containers[i] = edited
							}
							return nil
						}
					}(i, cont),
				})
		}
		containerMenuItems = append(containerMenuItems,
			[]*activekit.MenuItem{
				{
					Label: "Add new container",
					Action: func() error {
						logrus.Debugf("adding container")
						cont, ok := getContainer(container.Container{
							Container: model.Container{
								Name: namegen.Aster() + "-" + namegen.Color(),
								Limits: model.Resource{
									Memory: 256,
									CPU:    200,
								},
								Ports: []model.ContainerPort(nil),
							},
						})
						if ok {
							containers = append(containers, cont)
							fmt.Printf("Container %q added to list\n", cont.Name)
						}
						return nil
					},
				},
				{
					Label: "Delete container",
					Action: func() error {
						logrus.Debugf("deleting container")
						var deleteMenu []*activekit.MenuItem
						for i, name := range getContainersNamesList(containers) {
							deleteMenu = append(deleteMenu, &activekit.MenuItem{
								Label: name,
								Action: func(i int, name string) func() error {
									return func() error {
										yes, _ := activekit.Yes(fmt.Sprintf("Are you sure you want to delete the container %q?",
											name))
										if yes {
											containers = append(containers[:i], containers[i+1:]...)
										}
										return nil
									}
								}(i, name),
							})
						}
						deleteMenu = append(deleteMenu, &activekit.MenuItem{
							Label: "Return to previous menu",
						})
						(&activekit.Menu{
							Title: "Which container do you want to delete?",
							Items: deleteMenu,
						}).Run()
						return nil
					},
				},
				{
					Label: "Confirm",
					Action: func() error {
						exit = true
						ok = true
						return nil
					},
				},
				{
					Label: "Return to previous menu, drop all changes",
					Action: func() error {
						exit = true
						ok = false
						return nil
					},
				},
			}...)
		_, err := (&activekit.Menu{
			Items: containerMenuItems,
		}).Run()
		if err != nil {
			logrus.WithError(err).Errorf("")
			break
		}
	}
	if ok {
		return containers
	}
	return conts
}

func getName(defaultName string) string {
	for {
		name, _ := activekit.AskLine(fmt.Sprintf("Print deployment name (or hit Enter to use %q) > ", defaultName))
		if strings.TrimSpace(name) == "" {
			name = defaultName
		}
		if err := validation.ValidateLabel(name); err != nil {
			fmt.Printf("Invalid name %q. Try again\n", name)
			continue
		}
		return name
	}
}

func getReplicas(defaultReplicas int) int {
	for {
		replicasStr, _ := activekit.AskLine(fmt.Sprintf("Print number or replicas (%v, hit Enter to use %d) > ", ReplicasLimit, defaultReplicas))
		replicas := defaultReplicas
		if strings.TrimSpace(replicasStr) == "" {
			return defaultReplicas
		}
		if _, err := fmt.Sscan(replicasStr, &replicas); err != nil || !ReplicasLimit.Containing(replicas) {
			fmt.Printf("Expected number %v! Try again.\n", ReplicasLimit)
			continue
		}
		return replicas
	}
}

func getContainersNamesList(containers []container.Container) []string {
	names := make([]string, 0, len(containers))
	for _, cont := range containers {
		names = append(names, cont.Name)
	}
	return names
}
