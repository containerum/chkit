package deplactive

import (
	"fmt"
	"strings"

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
		if _, err := fmt.Sscan(replicasStr, &replicas); err != nil || replicas > 15 {
			fmt.Printf("Expected number 1..15! Try again.\n")
			continue
		}
		return replicas
	}
}

func getContainers() []container.Container {
	containers := []container.Container{}
	for {
		cont, ok := getContainer()
		if ok {
			containers = append(containers, cont)
			fmt.Printf("Container %q added to list\n", cont.Name)
		}
		yes, _ := activeToolkit.Yes("Add another container?")
		if !yes {
			return containers
		}
	}
}

func getContainer() (container.Container, bool) {
	con := container.Container{
		model.Container{
			Name:  namegen.Aster() + "-" + namegen.Color(),
			Image: "unknown (required)",
			Limits: model.Resource{
				Memory: "",
				CPU:    "",
			},
			Ports: []model.ContainerPort(nil),
		},
	}
	fmt.Printf("Ok, the hard part. Let's create a container\n")
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
		return resource.NewMilliQuantity(int64(mem), resource.BinarySI).String()
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
		if _, err := fmt.Sscanln(cpuStr, &cpu); err != nil {
			fmt.Printf("CPU must be number > 0. Try again.\n")
			continue
		}
		return resource.NewMilliQuantity(int64(1000*cpu), resource.BinarySI).String()
	}
}
