package deplactive

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/model/container"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/containerum/chkit/pkg/util/validation"
)

func getContainer(con container.Container) (container.Container, bool) {
	ok := true
	for exit := false; !exit; {
		(&activekit.Menu{
			Items: []*activekit.MenuItem{
				{
					Label: fmt.Sprintf("Set name         : %s", con.Name),
					Action: func() error {
						con.Name = getContainerName(con.Name)
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Set image        : %s",
						activekit.OrString(con.Image, "none (required)")),
					Action: func() error {
						con.Image = getContainerImage()
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Set memory limit : %d Mb", con.Limits.Memory),
					Action: func() error {
						con.Limits.Memory = getMemory(con.Limits.Memory)
						return nil
					},
				},
				{
					Label: fmt.Sprintf("Set CPU limit    : %d mCPU", con.Limits.CPU),
					Action: func() error {
						con.Limits.CPU = getCPU(con.Limits.CPU)
						return nil
					},
				},
				{
					Label: "Edit environment variables",
					Action: func() error {
						editContainerEnvironmentVars(&con)
						return nil
					},
				},
				{
					Label: "Confirm",
					Action: func() error {
						if err := ValidateContainer(con); err != nil {
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
					Label: "Return to previous menu",
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

func getMemory(oldValue uint) uint {
	for {
		memStr, _ := activekit.AskLine(fmt.Sprintf("Memory (Mb, %v) > ", MemLimit))
		memStr = strings.TrimSpace(memStr)
		var mem uint
		if memStr == "" {
			return oldValue
		}
		if _, err := fmt.Sscanln(memStr, &mem); err != nil || !MemLimit.Containing(int(mem)) {
			fmt.Printf("Memory must be interger number %v. Try again.\n", MemLimit)
			continue
		}
		return mem
	}
}

func getCPU(oldValue uint) uint {
	for {
		cpuStr, _ := activekit.AskLine(fmt.Sprintf("CPU (%v mCPU) > ", CPULimit))
		cpuStr = strings.TrimSpace(cpuStr)
		var cpu uint
		if cpuStr == "" {
			return oldValue
		}
		if _, err := fmt.Sscanln(cpuStr, &cpu); err != nil || !CPULimit.Containing(int(cpu)) {
			fmt.Printf("CPU must be number %v. Try again.\n", CPULimit)
			continue
		}
		return cpu
	}
}
