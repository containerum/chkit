package deplactive

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/util/activeToolkit"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/containerum/chkit/pkg/util/validation"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	ErrInvalidDeploymentName chkitErrors.Err = "invalid deployment name"
)

func getName(defaultName string) string {
	for {
		name, _ := activeToolkit.AskLine(fmt.Sprintf("Print deployment name (or hit Enter to use %q) > ", defaultName))
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

func getReplicas(defaultReplicas uint) uint {
	for {
		replicasStr, _ := activeToolkit.AskLine(fmt.Sprintf("Print number or replicas (1..15, hit Enter to user %d) > ", defaultReplicas))
		replicas := defaultReplicas
		if strings.TrimSpace(replicasStr) == "" {
			return defaultReplicas
		}
		if _, err := fmt.Sscan(replicasStr, &replicas); err != nil || replicas == 0 || replicas > 15 {
			fmt.Printf("Expected number 1..15! Try again.\n")
			continue
		}
		return replicas
	}
}

func getContainers(containers []container.Container) []container.Container {
	for {
		containerNames := make([]string, 0, len(containers))
		for _, cont := range containers {
			containerNames = append(containerNames,
				fmt.Sprintf("Edit container %q", cont.Name))
		}
		containersOptions := append(containerNames,
			"Add new container",
			"Delete container",
			"Exit")
		_, option, _ := activeToolkit.Options("What do you want?", false,
			containersOptions...)
		logrus.Debugf("option %d in %d %+v", option, len(containersOptions), containersOptions)
	containerMenu:
		switch option {
		case len(containersOptions) - 1: // Exit
			logrus.Debugf("exit container menu")
			return containers
		case len(containersOptions) - 2: // Delete container
			logrus.Debugf("delete container menu")
			_, option, _ := activeToolkit.Options("Which container do you want to delete?", false,
				append(containerNames, "Exit")...)
			switch option {
			case len(containerNames): // exit
				break containerMenu
			default:
				if len(containers) > 0 {
					logrus.Debugf("delete container %q", containerNames[option])
					containers = append(containers[:option], containers[option+1:]...)
				}
			}
		case len(containersOptions) - 3: // Add container
			logrus.Debugf("add container")
			cont, ok := getContainer(container.Container{
				model.Container{
					Name:  namegen.Aster() + "-" + namegen.Color(),
					Image: "unknown (required)",
					Limits: model.Resource{
						Memory: "",
						CPU:    "",
					},
					Ports: []model.ContainerPort(nil),
				},
			})
			if ok {
				containers = append(containers, cont)
				fmt.Printf("Container %q added to list\n", cont.Name)
			}
		default: // edit container
			logrus.Debugf("editing container %q", containerNames[option])
			edited, ok := getContainer(containers[option])
			if ok {
				containers[option] = edited
			}
		}
	}
}

func getContainer(con container.Container) (container.Container, bool) {
	for {
		_, option, _ := activeToolkit.Options("Choose option: ", false,
			fmt.Sprintf("Set name         : %s", con.Name),
			fmt.Sprintf("Set image        : %s",
				activeToolkit.OrString(con.Image, "none (required)")),
			fmt.Sprintf("Set memory limit : %s",
				activeToolkit.OrString(con.Limits.Memory, "none (required)")),
			fmt.Sprintf("Set CPU limit    : %s",
				activeToolkit.OrString(con.Limits.CPU, "none (requied)")),
			"Confirm",
			"Exit")
		switch option {
		case 0:
			con.Name = getContainerName(con.Name)
		case 1:
			con.Image = getContainerImage()
		case 2:
			con.Limits.Memory = getMemory()
		case 3:
			con.Limits.CPU = getCPU()
		case 4:
			if err := validateContainer(con); err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			return con, true
		default:
			return con, false
		}
	}
}

func getContainerName(defaultName string) string {
	for {
		name, _ := activeToolkit.AskLine(fmt.Sprintf("Type container name (press Enter to use %q) > ", defaultName))
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
		image, _ := activeToolkit.AskLine("> ")
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

func getLimits() model.Resource {
	return model.Resource{
		CPU:    getCPU(),
		Memory: getMemory(),
	}
}

func getMemory() string {
	for {
		memStr, _ := activeToolkit.AskLine("Memory (Mb) > ")
		memStr = strings.TrimSpace(memStr)
		var mem uint32
		if memStr == "" {
			return ""
		}
		if _, err := fmt.Sscanln(memStr, &mem); err != nil {
			fmt.Printf("Memory must be interger number > 0. Try again.\n")
			continue
		}
		return resource.NewQuantity(int64(mem*(1<<20)), resource.BinarySI).String()
	}
}

func getCPU() string {
	for {
		cpuStr, _ := activeToolkit.AskLine("CPU (0.6 of CPU for example) > ")
		cpuStr = strings.TrimSpace(cpuStr)
		var cpu float32
		if cpuStr == "" {
			return ""
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
