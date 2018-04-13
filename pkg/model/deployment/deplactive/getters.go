package deplactive

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/containerum/chkit/pkg/util/validation"
	"k8s.io/apimachinery/pkg/api/resource"
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
					Name: fmt.Sprintf("Edit container %q", cont.Name),
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
					Name: "Add new container",
					Action: func() error {
						logrus.Debugf("adding container")
						cont, ok := getContainer(container.Container{
							Container: model.Container{
								Name:  namegen.Aster() + "-" + namegen.Color(),
								Image: "unknown (required)",
								Limits: model.Resource{
									Memory: "256Mi",
									CPU:    "200m",
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
					Name: "Delete container",
					Action: func() error {
						logrus.Debugf("deleting container")
						var deleteMenu []*activekit.MenuItem
						for i, name := range getContainersNamesList(containers) {
							deleteMenu = append(deleteMenu, &activekit.MenuItem{
								Name: name,
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
							Name: "Return to previous menu",
						})
						(&activekit.Menu{
							Title: "Which container do you want to delete?",
							Items: deleteMenu,
						}).Run()
						return nil
					},
				},
				{
					Name: "Confirm",
					Action: func() error {
						exit = true
						ok = true
						return nil
					},
				},
				{
					Name: "Return to previous menu, drop all changes",
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

func getContainer(con container.Container) (container.Container, bool) {
	ok := true
	for exit := false; !exit; {
		(&activekit.Menu{
			Items: []*activekit.MenuItem{
				{
					Name: fmt.Sprintf("Set name         : %s", con.Name),
					Action: func() error {
						con.Name = getContainerName(con.Name)
						return nil
					},
				},
				{
					Name: fmt.Sprintf("Set image        : %s",
						activekit.OrString(con.Image, "none (required)")),
					Action: func() error {
						con.Image = getContainerImage()
						return nil
					},
				},
				{
					Name: fmt.Sprintf("Set memory limit : %s",
						activekit.OrString(con.Limits.Memory, "none (required)")),
					Action: func() error {
						con.Limits.Memory = getMemory(con.Limits.Memory)
						return nil
					},
				},
				{
					Name: fmt.Sprintf("Set CPU limit    : %s",
						activekit.OrString(con.Limits.CPU, "none (requied)")),
					Action: func() error {
						con.Limits.CPU = getCPU(con.Limits.CPU)
						return nil
					},
				},
				{
					Name: "Confirm",
					Action: func() error {
						if err := validateContainer(con); err != nil {
							errText := err.Error()
							attention := strings.Repeat("!", text.Width(errText))
							fmt.Printf("%s\n%v\n%s\n", attention, errText, attention)
							return nil
						}
						exit = true
						return nil
					},
				},
				{
					Name: "Return to previous menu",
					Action: func() error {
						ok = false
						exit = true
						return nil
					},
				},
			},
		}).Run()
	}
	return con, ok
}

func getContainerName(defaultName string) string {
	for {
		name, _ := activekit.AskLine(fmt.Sprintf("Type container name (press Enter to use %q) > ", defaultName))
		name = strings.TrimSpace(name)
		if name == "" {
			name = defaultName
		}
		if validation.ValidateContainerName(name) != nil {
			fmt.Printf("Invalid name :( Try again.\n")
			continue
		}
		return name
	}
}

func getContainerImage() string {
	fmt.Printf("Which image do you want to use?\n")
	for {
		image, _ := activekit.AskLine("> ")
		image = strings.TrimSpace(image)
		if image == "" {
			return ""
		}
		if validation.ValidateImageName(image) != nil {
			fmt.Printf("Invalid image name :( Try again.\n")
			continue
		}
		return image
	}
}

func getMemory(oldValue string) string {
	for {
		memStr, _ := activekit.AskLine("Memory (Mb, 1..16000) > ")
		memStr = strings.TrimSpace(memStr)
		var mem uint32
		if memStr == "" {
			return oldValue
		}
		if _, err := fmt.Sscanln(memStr, &mem); err != nil || mem < 1 || mem > 16000 {
			fmt.Printf("Memory must be interger number 0..16000. Try again.\n")
			continue
		}
		return resource.NewQuantity(int64(mem*(1<<20)), resource.BinarySI).String()
	}
}

func getCPU(oldValue string) string {
	for {
		cpuStr, _ := activekit.AskLine("CPU (0.6 of CPU for example, 0.001..12) > ")
		cpuStr = strings.TrimSpace(cpuStr)
		var cpu float32
		if cpuStr == "" {
			return oldValue
		}
		if _, err := fmt.Sscanln(cpuStr, &cpu); err != nil || cpu <= 0 {
			fmt.Printf("CPU must be number > 0. Try again.\n")
			continue
		}
		cpuQ := resource.NewScaledQuantity(int64(1000*cpu), resource.Milli)
		cpuQ.Format = resource.BinarySI
		return cpuQ.String()
	}
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
		replicasStr, _ := activekit.AskLine(fmt.Sprintf("Print number or replicas (1..15, hit Enter to use %d) > ", defaultReplicas))
		replicas := defaultReplicas
		if strings.TrimSpace(replicasStr) == "" {
			return defaultReplicas
		}
		if _, err := fmt.Sscan(replicasStr, &replicas); err != nil || replicas < 0 || replicas > 15 {
			fmt.Printf("Expected number 1..15! Try again.\n")
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
